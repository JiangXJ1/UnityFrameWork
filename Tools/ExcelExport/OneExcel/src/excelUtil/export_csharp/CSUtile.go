package export_csharp

import (
	"OneExcel/src"
	"strings"
)

func NeedExportCSorBinary(sheetData src.SheetData) bool {
	containBin := false
	for i := range sheetData.FieldExportTypeList {
		if i > 0 {
			exportType := sheetData.FieldExportTypeList[i]
			if strings.Contains(exportType, "C") {
				containBin = true
				break
			}
		}
	}
	return containBin
}

func NeedExportCSV(sheetData src.SheetData) bool {
	containBin := false
	for i := range sheetData.FieldExportTypeList {
		if i > 0 {
			exportType := sheetData.FieldExportTypeList[i]
			if strings.Contains(exportType, "S") {
				containBin = true
				break
			}
		}
	}
	return containBin
}
