package binary

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/export_csharp"
	"OneExcel/src/excelUtil/tool"
	"strings"
)

func WriteGlobalBinary(sheetData src.SheetData) {
	fileStresam := export_csharp.NewStream()
	rowStream := export_csharp.NewStream()

	contentList := tool.GetGlobalContent(sheetData)
	for _, cellValues := range contentList {
		exportable := tool.CanExportCS(cellValues.ExportType, true)
		if !exportable {
			continue
		}
		//key := strings.TrimSpace(cellValues.Key)
		contentType := strings.TrimSpace(cellValues.ContentType)
		content := strings.TrimSpace(cellValues.Content)
		//des := cellValues.Des

		rowStream.WriteNodeValue(contentType, content)
	}

	fileStresam.WriteInt32(int32(0))
	fileStresam.WriteInt32(int32(rowStream.Len()))
	fileStresam.WriteRawBytes(rowStream.Buffer().Bytes())

	fileStresam.WriteFile(src.GetBinFilePathName(strings.ReplaceAll(sheetData.SheetName, "@", "")))
}
