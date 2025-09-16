package csv

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/export_csharp"
	"OneExcel/src/fileUtil"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func WriteCSV(sheetData src.SheetData) {
	containBin := export_csharp.NeedExportCSV(sheetData)
	if containBin == false {
		containBin := export_csharp.NeedExportCSV(sheetData)
		if containBin == false {
			fmt.Println(fmt.Sprintf("%v---%v:不会导出服务器csv配置", sheetData.FileName, sheetData.SheetName))
			return
		}
		return
	}
	if sheetData.SheetName == "Schema_Global" || strings.HasPrefix(sheetData.SheetName, "Schema_@") {
		WriteGlobalCSV(sheetData)
	} else if strings.Contains(sheetData.SheetName, "Language_") {
		WriteLanguageCSV(sheetData)
	} else {
		WriteTableCSV(sheetData)
	}
}

func WriteTableCSV(sheetData src.SheetData) {
	fileUtil.CreateDir(src.OUTPUT_DIR_CSV)
	filePath := src.GetCSVFilePathName(strings.ReplaceAll(sheetData.SheetName, "@", ""))
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	write := bufio.NewWriter(file)

	fieldNameStr := ""
	fieldTypeStr := ""
	splitStr := ""
	count := len(sheetData.FieldExportTypeList)
	for j := range sheetData.FieldExportTypeList {
		exportType := sheetData.FieldExportTypeList[j]
		if !strings.Contains(exportType, "S") {
			continue
		}
		if j == count-1 {
			splitStr = ""
		} else {
			splitStr = ","
		}
		fieldNameStr = fieldNameStr + sheetData.FieldNameList[j] + splitStr
		fieldTypeStr = fieldTypeStr + sheetData.FieldTypeList[j] + splitStr
	}
	write.WriteString(fieldNameStr + "\n")
	write.WriteString(fieldTypeStr + "\n")

	str_line := ""
	for i := range sheetData.AllRowValueList {
		str_line = ""
		//行
		cellValues := sheetData.AllRowValueList[i]
		for j := range cellValues {

			if j >= count {
				continue
			}
			//列
			cellValue := cellValues[j]
			exportType := sheetData.FieldExportTypeList[j]
			if !strings.Contains(exportType, "S") {
				continue
			}

			cellValue = strings.TrimSpace(cellValue)
			if len(cellValue) == 0 {
				defaultValue := strings.TrimSpace(sheetData.DefaultValueList[j])
				cellValue = defaultValue
			}
			if j == count-1 {
				str_line = str_line + cellValue
			} else {
				str_line = str_line + cellValue + ","
			}
		}
		if str_line != "" {
			write.WriteString(str_line + "\n")
		}
	}
	write.Flush()
}

func init() {

}
