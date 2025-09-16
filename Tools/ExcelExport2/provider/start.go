package provider

import (
	"context"
	"demo-go-excel/exports/cs/export"
	"demo-go-excel/global"
	"demo-go-excel/util"
	"demo-go-excel/util/fileUtil"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	// 创建输出目录
	fileUtil.CreateDirectoryIfNotExists(global.OUTPUT_DIR_BIN)
	fileUtil.CreateDirectoryIfNotExists(global.OUTPUT_DIR_CS)
	fileUtil.CreateDirectoryIfNotExists(global.OUTPUT_DIR_LUA)
	fileUtil.CreateDirectoryIfNotExists(global.OUTPUT_DIR_CSV)
}

func Start() {
	startTime := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	runInfo := NewRunInfo(ctx)
	defer func() {
		cancel()
		close(runInfo.writer)
		close(runInfo.errChan)
	}()

	// 启动错误监听
	go func() {
		for err := range runInfo.errChan {
			fmt.Print("\n" + err.Error())
			runInfo.SetError()
			cancel() // 发生错误时取消上下文
		}
	}()

	files, _ := os.ReadDir(global.EXCEL_DIR)

	// 读取md5文件
	//md5Map, _ := util.LoadHashesFromFile(global.MD5_FILE)

	//启动线程池
	goroutinePool := NewGoroutinePool(128, 500)
	goroutinePool.Start()

	global.MaxFileCount = int32(len(files))
	waitGroup := sync.WaitGroup{}
	for _, file := range files {
		waitGroup.Add(1)
		go func(file os.DirEntry, s *sync.WaitGroup) {
			defer func() {
				s.Done()
				atomic.AddInt32(&global.CurFileCount, 1)
				// 只在没有错误时显示进度条
				if !runInfo.HasError() {
					fmt.Printf("\r[%-50s] 文件：%d/%d  页签：%d",
						getProgressBar(global.CurFileCount, global.MaxFileCount),
						global.CurFileCount,
						global.MaxFileCount,
						global.CurTaskCount)
				}
			}()
			fileName := file.Name()
			if file.IsDir() || !util.IsExcelFile(fileName) || strings.Contains(fileName, "~$") {
				return
			}
			filePath := global.EXCEL_DIR + fileName
			// 计算md5值
			//md5Value, err := util.CalculateFileHash(filePath)

			isDirty := true
			//if err == nil && md5Map[fileName] == md5Value {
			//	isDirty = false
			//	fmt.Printf(">>>当前文件未修改:[%v] \n", fileName)
			//}
			//md5Map[fileName] = md5Value
			//startTime1 := time.Now()
			f, err := xlsx.OpenFile(filePath)
			if err != nil {
				fmt.Printf(">>>打开文件失败:[%v] \n", filePath)
				return
			}

			for _, sheet := range f.Sheets {
				if strings.HasPrefix(sheet.Name, "#") {
					continue
				}
				//添加任务
				goroutinePool.Submit(&Task{
					runInfo,
					sheet,
					isDirty,
					fileName,
				})
			}
			//fmt.Printf("打开文件耗时:%s  %s\n", fileName, time.Now().Sub(startTime1))
		}(file, &waitGroup)
	}

	//等待添加任务完成
	waitGroup.Wait()
	//等待所有线程完成
	goroutinePool.WaitComplete()

	if runInfo.HasError() {
		fmt.Println("\n导出过程中出现错误，部分表格可能未完成导出!")
	} else {
		fmt.Println("\n所有任务执行完成!")
		CreateSchemaCreator()
	}

	// 更新md5文件
	//err := util.WriteHashesToFile(global.MD5_FILE, md5Map)
	//if err != nil {
	//	fmt.Printf("更新md5文件失败:%v\n", err)
	//}

	fmt.Printf("耗时:%.1f秒\n", time.Now().Sub(startTime).Seconds())

	if runtime.GOOS == "windows" {
		fmt.Println("按任意键退出...")
		fmt.Scanln() // 等待用户输入，防止窗口关闭
	}
}

func CreateSchemaCreator() {
	outputDir := global.OUTPUT_DIR_CS
	filePath := filepath.Join(outputDir, "SchemaCreator.cs")

	// 读取目录下所有文件
	files, err := os.ReadDir(outputDir)
	if err != nil {
		panic("读取目录失败: " + err.Error())
	}

	// 过滤并收集有效的文件名
	var validFiles []string
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".cs" || file.Name() == "SchemaCreator.cs" {
			continue
		}
		validFiles = append(validFiles, file.Name())
	}

	// 对文件名进行排序
	sort.Strings(validFiles)

	var enumContent strings.Builder
	var switchContent strings.Builder

	for index, fileName := range validFiles {
		// 获取文件名（不含后缀）
		baseName := strings.TrimSuffix(fileName, ".cs")

		// 构建枚举内容
		enumContent.WriteString("\t\t")
		enumContent.WriteString(baseName)
		if index < len(validFiles)-1 {
			enumContent.WriteString(",\n")
		} else {
			enumContent.WriteString(",")
		}

		// 构建switch case内容
		switchContent.WriteString("\t\t\t\tcase SchemaID.")
		switchContent.WriteString(baseName)
		switchContent.WriteString(": return new Schema.")
		switchContent.WriteString(baseName)
		if index < len(validFiles)-1 {
			switchContent.WriteString("();\n")
		} else {
			switchContent.WriteString("();")
		}
	}

	// 使用模板生成内容
	param := &export.SchemaCreatorTemplateParam{
		EnumContent: enumContent.String(),
		Content:     switchContent.String(),
	}

	content := param.GenerateCsharpTemplate()

	// 写入文件
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		panic("写入文件失败: " + err.Error())
	}
}

// getProgressBar 返回一个表示进度的字符串
func getProgressBar(current, total int32) string {
	progress := float64(current) / float64(total)
	barLength := 50
	filled := int(progress * float64(barLength))
	bar := ""
	for i := 0; i < filled; i++ {
		bar += "="
	}
	for i := filled; i < barLength; i++ {
		bar += " "
	}
	return bar
}
