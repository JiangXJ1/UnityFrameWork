package provider

import (
	"bufio"
	"demo-go-excel/exports/cs/export"
	"demo-go-excel/global"
	"demo-go-excel/util/common"
	"demo-go-excel/util/fileUtil"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
)

type LanguageType struct {
	Column      int
	BinFileName string
	ContentList []LanguageContent
}

type LanguageContent struct {
	Key     string
	Content string
}

type SheetInfoLanguage struct {
	FileName string
	Sheet    *xlsx.Sheet
	Types    []LanguageType
}

func (s *SheetInfoLanguage) Read(task *Task) {
	s.FileName = task.FileName
	if len(task.Sheet.Rows) < global.MIN_ROW_COUNT {
		return
	}
	s.Sheet = task.Sheet
	//解析字段信息和内容
	s.GetFieldInfo()
}

// 解析字段信息
func (s *SheetInfoLanguage) GetFieldInfo() {
	firstRows := s.Sheet.Rows[0].Cells
	s.Types = make([]LanguageType, 0, len(firstRows)-3)
	rowLen := len(s.Sheet.Rows)
	columnCount := rowLen - 5
	for i, cell := range firstRows {
		if i < 3 {
			continue
		}
		if len(strings.TrimSpace(cell.Value)) > 0 {
			fieldName := s.Sheet.Rows[2].Cells[i].String()
			if len(fieldName) == 0 {
				continue
			}
			binFileName := strings.ReplaceAll(fieldName, "*", "")
			s.Types = append(s.Types, LanguageType{Column: i, BinFileName: binFileName, ContentList: make([]LanguageContent, 0, columnCount)})
		}
	}
	typeCount := len(s.Types)
	if typeCount == 0 {
		return
	}

	names := make(map[string]bool)
	for i := 5; i < rowLen; i++ {
		cells := s.Sheet.Rows[i].Cells
		key := cells[1].Value
		if len(cells) < 4 || len(strings.TrimSpace(key)) == 0 || strings.HasPrefix(strings.TrimSpace(cells[0].Value), "#") {
			continue
		}

		if names[key] {
			fmt.Println("字段名重复：" + key)
			panic("字段名重复：" + key)
		}
		names[key] = true
		for i := range s.Types {
			lanContent := cells[s.Types[i].Column].Value
			s.Types[i].ContentList = append(s.Types[i].ContentList, LanguageContent{Key: key, Content: lanContent})
		}
	}
}

// Export 导出
func (s *SheetInfoLanguage) Export() {
	if len(s.Types) == 0 {
		return
	}

	s.ExportCsharp()
	s.ExportBin()
}

// ExportCsharp 导出C#代码
func (s *SheetInfoLanguage) ExportCsharp() {
	name := global.LANGUAGE_NAME
	filePath := fileUtil.GetCSFilePathName(name)
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(file)

	//模版参数
	CsharpGlobalTemplateParam := export.CsharpLanguageTemplateParam{
		ClassName:  fmt.Sprintf("Schema_%s", strings.ReplaceAll(name, "@", "")),
		DataDefine: "public Dictionary<string, string> language = new Dictionary<string, string>();",
	}
	//生成代码
	content := CsharpGlobalTemplateParam.GenerateCsharpTemplate()
	//写入文件
	write := bufio.NewWriter(file)
	_, err = write.WriteString(content)
	if err != nil {
		return
	}
	err = write.Flush()
	if err != nil {
		return
	}
}

// ExportBin 导出二进制文件
func (s *SheetInfoLanguage) ExportBin() {
	for _, info := range s.Types {
		fileStream := common.NewStream()
		rowStream := common.NewStream()
		rowStream.WriteInt32(int32(len(info.ContentList)))
		for _, item := range info.ContentList {
			rowStream.WriteNodeValue("string", item.Key)
			rowStream.WriteNodeValue("string", item.Content)
		}
		fileStream.WriteInt32(int32(0))
		fileStream.WriteInt32(int32(rowStream.Len()))
		fileStream.WriteRawBytes(rowStream.Buffer().Bytes())

		err := fileStream.WriteFile(fileUtil.GetBinFilePathName(info.BinFileName))
		if err != nil {
			fmt.Printf("导出二进制文件失败：%s\n", err.Error())
		}
	}
}
