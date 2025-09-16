package cs

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/tool"
	"OneExcel/src/fileUtil"
	"bufio"
	"log"
	"os"
	"strings"
)

func WriteGlobalCS(sheetData src.SheetData) {
	fileUtil.CreateDir(src.OUTPUT_DIR_CS)
	filePath := src.GetCSFilePathName(strings.ReplaceAll(sheetData.SheetName, "@", ""))
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	str_define := ""
	str_deserialize := ""

	contentList := tool.GetGlobalContent(sheetData)
	rowLength := len(contentList)
	for i, cellValues := range contentList {
		exportable := tool.CanExportCS(cellValues.ExportType, true)
		if !exportable {
			continue
		}
		key := strings.TrimSpace(cellValues.Key)
		contentType := strings.TrimSpace(cellValues.ContentType)
		//content := strings.TrimSpace(cellValues.Content)
		des := cellValues.Des
		des = strings.ReplaceAll(des, "\n", "。")

		content_define := GetTypeStr(contentType, key)
		content_define = strings.ReplaceAll(content_define, "//*", "//"+des)
		content_deser := GetDeserializeStr(contentType, key)

		if i == rowLength-1 {
			str_define = str_define + content_define
			str_deserialize = str_deserialize + content_deser
		} else {
			str_define = str_define + content_define + "\n\t\t"
			str_deserialize = str_deserialize + content_deser + "\n\n\t\t\t"
		}
	}

	str_dataDefine := str_define
	str_dataDeserialize := str_deserialize

	content := CsharpGlobalTemplate
	content = strings.ReplaceAll(content, "#KEY_NAME#", sheetData.FieldNameList[0])
	content = strings.ReplaceAll(content, "#DATA_DEFINE#", str_dataDefine)
	content = strings.ReplaceAll(content, "#DESERIALIZE#", str_dataDeserialize)
	content = strings.ReplaceAll(content, "#CLASS_NAME#", strings.ReplaceAll(sheetData.SheetName, "@", ""))

	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
}
