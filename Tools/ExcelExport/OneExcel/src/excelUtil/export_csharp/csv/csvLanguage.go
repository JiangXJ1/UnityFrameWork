package csv

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/export_csharp"
	"OneExcel/src/excelUtil/tool"
)

func WriteLanguageCSV(sheetData src.SheetData) {
	fileStresam := export_csharp.NewStream()
	rowStream := export_csharp.NewStream()

	contentList := tool.GetLanguageData(sheetData)
	for _, value := range contentList {
		exportType := value.ExportType
		exportable := tool.CanExportBin(exportType, true)
		if !exportable {
			continue
		}
		//key := value[1]
		//des := value[2]
		content := value.Content
		rowStream.WriteNodeValue("string", content)
	}

	fileStresam.WriteInt32(int32(0))
	fileStresam.WriteInt32(int32(rowStream.Len()))
	fileStresam.WriteRawBytes(rowStream.Buffer().Bytes())

	fileStresam.WriteFile(src.GetBinFilePathName(sheetData.SheetName))
}
