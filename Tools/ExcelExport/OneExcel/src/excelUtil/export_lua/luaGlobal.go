package export_lua

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/tool"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func WriteGlobal(sheetData src.SheetData, write *bufio.Writer, isClient bool) {
	contentList := tool.GetGlobalContent(sheetData)
	//rowLength := len(contentList)

	write.WriteString(fmt.Sprintf("---@class %v\n", strings.ReplaceAll(sheetData.SheetName, "@", "")))
	for _, cellValues := range contentList {
		key := strings.TrimSpace(cellValues.Key)
		contentType := strings.TrimSpace(cellValues.ContentType)
		if len(key) > 0 {
			write.WriteString(fmt.Sprintf("---@field %v %v  @%v \n", key, contentType, strings.ReplaceAll(cellValues.Des, "\n", " ")))
		}
	}
	//Write value
	if isClient {
		write.WriteString("local Data =")
	} else {
		write.WriteString("local ")
		write.WriteString(strings.ReplaceAll(sheetData.SheetName, "@", ""))
		write.WriteString(" =")
	}

	write.WriteString(" \n{\n")

	for i, cellValues := range contentList {
		exportable := tool.CanExportLua(cellValues.ExportType, isClient)
		if !exportable {
			continue
		}
		key := strings.TrimSpace(cellValues.Key)
		contentType := strings.TrimSpace(cellValues.ContentType)
		content := strings.TrimSpace(cellValues.Content)
		//des := cellValues.Des

		write.WriteString("\t")
		write.WriteString(key)
		write.WriteString(" = ")

		fixCellValue, err := GetLuaCellValue(content, contentType)
		if err != nil {
			tip := sheetData.SheetName + " 行:" + strconv.Itoa(i)
			panic(err.Error() + tip)
		}
		write.WriteString(fixCellValue)

		//if i != rowLength-1 {
		//	write.WriteString(",\t--" + des)
		//} else {
		//	write.WriteString("\t--" + des)
		//}

		write.WriteString(",\n")
	}
}
