package export_lua

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/tool"
	"bufio"
)

func WriteLanguage(sheetData src.SheetData, write *bufio.Writer, isClient bool) {
	//Write value
	if isClient {
		write.WriteString("local Data =")
	} else {
		write.WriteString("local ")
		write.WriteString(sheetData.SheetName)
		write.WriteString(" =")
	}
	write.WriteString(" \n{\n")

	contentList := tool.GetLanguageData(sheetData)
	rowLength := len(contentList)

	for i, value := range contentList {
		exportable := tool.CanExportLua(value.ExportType, isClient)
		if !exportable {
			continue
		}

		key := value.Key
		des := value.Des
		content := value.Content

		write.WriteString("\t")
		write.WriteString(key)
		write.WriteString(" = ")
		write.WriteString("\"" + content + "\"")

		if i != rowLength-1 {
			write.WriteString(",\t--" + des)
			write.WriteString("\n")
		} else {
			write.WriteString("\t--" + des)
			write.WriteString("\n")
		}
	}
}
