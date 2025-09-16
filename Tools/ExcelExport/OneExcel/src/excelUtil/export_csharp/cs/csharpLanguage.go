package cs

import (
	"OneExcel/src"
	"OneExcel/src/fileUtil"
	"bufio"
	"log"
	"os"
	"strings"
)

func WriteLanguageCS(sheetData src.SheetData) {
	fileUtil.CreateDir(src.OUTPUT_DIR_CS)
	filePath := src.GetCSFilePathName("Schema_Language")
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	//str_define := ""
	//str_deserialize := ""

	//contentList := tool.GetLanguageData(sheetData)
	//rowLength := len(contentList)
	//str := `ins.language.Add("#NAME#", reader.ReadString());`
	//for i, value := range contentList {
	//	exportable := tool.CanExportCS(value.ExportType, true)
	//	if !exportable {
	//		continue
	//	}
	//
	//	key := value.Key
	//	//des := value.Des
	//	//content := value.Content
	//
	//	//content_define := GetTypeStr("string", key)
	//	//content_define = strings.ReplaceAll(content_define, "//*", "//"+des)
	//	content_deser := strings.ReplaceAll(str, "#NAME#", key)
	//
	//	if i == rowLength-1 {
	//		//str_define = str_define + content_define
	//		str_deserialize = str_deserialize + content_deser
	//	} else {
	//		//str_define = str_define + content_define + "\n\t\t"
	//		str_deserialize = str_deserialize + content_deser + "\n\t\t\t"
	//	}
	//}

	str_dataDefine := "public Dictionary<string, string> language = new Dictionary<string, string>();"
	//str_dataDeserialize := str_deserialize

	content := CsharpLanguageTemplate
	content = strings.ReplaceAll(content, "#KEY_NAME#", sheetData.FieldNameList[0])
	content = strings.ReplaceAll(content, "#DATA_DEFINE#", str_dataDefine)
	//content = strings.ReplaceAll(content, "#DESERIALIZE#", str_dataDeserialize)
	content = strings.ReplaceAll(content, "#CLASS_NAME#", "Schema_Language")

	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
}
