package csv

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/tool"
	"OneExcel/src/fileUtil"
	"bufio"
	"log"
	"os"
	"strings"
)

func WriteGlobalCSV(sheetData src.SheetData) {
	fileUtil.CreateDir(src.OUTPUT_DIR_CSV)
	filePath := src.GetCSVFilePathName(strings.ReplaceAll(sheetData.SheetName, "@", ""))
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	write := bufio.NewWriter(file)

	write.WriteString("KEY,Type,Content\n")
	str_line := ""

	contentList := tool.GetGlobalContent(sheetData)
	rowLength := len(contentList)
	for i, cellValues := range contentList {
		exportable := tool.CanExportSCSV(cellValues.ExportType)
		if !exportable {
			continue
		}
		key := strings.TrimSpace(cellValues.Key)
		contentType := strings.TrimSpace(cellValues.ContentType)
		if i == rowLength-1 {
			str_line = key + "," + contentType + "," + cellValues.Content
		} else {
			str_line = key + "," + contentType + "," + cellValues.Content + "\n"
		}
		write.WriteString(str_line)
	}
	write.Flush()
}
