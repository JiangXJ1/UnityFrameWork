package excelUtil

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/export_csharp/binary"
	"OneExcel/src/excelUtil/export_csharp/cs"
	"OneExcel/src/excelUtil/export_csharp/csv"
	"OneExcel/src/excelUtil/export_lua"
	"OneExcel/src/excelUtil/xlsx"
	"OneExcel/src/fileUtil"
)

func GetData(excelFilePath string, fileName string) []src.SheetData {
	//sheetDataList := excelize.GetData(config.EXCEL_ONE_PATH)
	return xlsx.GetData(excelFilePath, fileName)
}

func WriteData(sheetData src.SheetData) {
	export_lua.Write(sheetData)
	binary.WriteBinary(sheetData)
	cs.WriteTableCS(sheetData)
	csv.WriteCSV(sheetData)
}

func IsExcelClosed() bool {
	return fileUtil.IsExcelClosed()
}

func GetAllData() []src.SheetData {
	var results []src.SheetData
	fileNames := fileUtil.GetAllExcelName()
	for _, fileName := range fileNames {
		filePath := src.GetExcelFilePathName(fileName)
		excelSheetData := GetData(filePath, fileName)
		for _, excelData := range excelSheetData {
			results = append(results, excelData)
		}
	}
	return results
}
