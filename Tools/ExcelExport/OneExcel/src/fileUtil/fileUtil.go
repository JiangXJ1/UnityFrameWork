package fileUtil

import (
	"OneExcel/src"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func DeleteDir(dirPath string) {
	os.RemoveAll(dirPath)
}

// 按过滤列表删除文件
func DeleteFiles(dirPath string, filterList map[string]bool) {
	files, _ := os.ReadDir(dirPath)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileName := f.Name()
		if filterList == nil || len(filterList) == 0 {
			DeleteFile(dirPath + fileName)
			continue
		}
		if filterList[fileName] {
			DeleteFile(dirPath + fileName)
		}
	}
}

func DeleteFile(fileName string) {
	var err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func CreateDir(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func FileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateFile(filePath string) {
	if FileExist(filePath) {
		DeleteFile(filePath)
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func IsExcelClosed() bool {
	hasExcelOpen := false
	files, _ := ioutil.ReadDir(src.EXCEL_DIR)
	for _, f := range files {
		fileName := f.Name()
		if strings.HasPrefix(fileName, "~") {
			hasExcelOpen = true
			log.Fatalf("Please Close Excel: %s", fileName)
		}
	}
	if hasExcelOpen {
		return false
	}
	return true
}

func CalculateFileHash(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	_, err = hasher.Write(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// 读取文件并将内容存入map
func LoadHashesFromFile(filePath string) (map[string]string, error) {
	hashMap := make(map[string]string)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "===")
		if len(parts) == 2 {
			hashMap[parts[0]] = parts[1]
		} else {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件MD5失败")
		return hashMap, err
	}

	return hashMap, nil
}

// 将map内容写回到文件
func WriteHashesToFile(filePath string, hashMap map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for name, hash := range hashMap {
		_, err := writer.WriteString(fmt.Sprintf("%s===%s\n", name, hash))
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func GetAllExcelName() []string {
	var results []string
	files, _ := ioutil.ReadDir(src.EXCEL_DIR)

	// 读取md5文件
	md5Map, _ := LoadHashesFromFile(src.EXCEL_MD5_PATH)

	for _, f := range files {
		fileName := f.Name()
		//过滤掉临时文件
		if strings.Contains(fileName, "~$") {
			fmt.Println("------------------------------正在打开的表，请注意保存: " + fileName)
			continue
		}

		// 计算md5值
		md5Value, err := CalculateFileHash(src.EXCEL_DIR + fileName)
		// 如果md5值相同，则跳过
		if err == nil && md5Map[fileName] == md5Value {
			//TODO 需要先解决缓存问题（目前是每次导出前将所有配置都删除了）
			//fmt.Printf(">>>当前表未修改:[%v]:[%v] \n", fileName, md5Value)
			//continue
		}

		md5Map[fileName] = md5Value

		//统计数据
		if strings.HasSuffix(fileName, ".xlsx") || strings.HasSuffix(fileName, ".xls") || strings.HasSuffix(fileName, ".xlsm") {
			fmt.Println(fileName)
			results = append(results, fileName)
		}
	}
	// 写入md5文件
	WriteHashesToFile(src.EXCEL_MD5_PATH, md5Map)

	return results
}
