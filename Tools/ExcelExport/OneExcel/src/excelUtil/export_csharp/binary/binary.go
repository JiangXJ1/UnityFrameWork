package binary

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/export_csharp"
	"fmt"
	"strings"
)

func WriteBinary(sheetData src.SheetData) {
	containBin := export_csharp.NeedExportCSorBinary(sheetData)
	if containBin == false {
		containBin := export_csharp.NeedExportCSorBinary(sheetData)
		if containBin == false {
			fmt.Println(fmt.Sprintf("%v---%v:不会导出BIN配置", sheetData.FileName, sheetData.SheetName))
			return
		}
		return
	}
	if sheetData.SheetName == "Schema_Global" || strings.HasPrefix(sheetData.SheetName, "Schema_@") {
		WriteGlobalBinary(sheetData)
	} else if strings.Contains(sheetData.SheetName, "Language_") {
		WriteLanguageBinary(sheetData)
	} else {
		WriteTableBinary(sheetData)
	}
}

func WriteTableBinary(sheetData src.SheetData) {
	fileStresam := export_csharp.NewStream()
	exportTypeCount := len(sheetData.FieldExportTypeList)
	for i := range sheetData.AllRowValueList {
		rowStream := export_csharp.NewStream()
		//行
		cellValues := sheetData.AllRowValueList[i]
		for j := range cellValues {

			if j >= exportTypeCount {
				continue
			}
			//列
			cellValue := cellValues[j]
			_cellType := sheetData.FieldTypeList[j]

			exportType := sheetData.FieldExportTypeList[j]
			if !strings.Contains(exportType, "C") {
				continue
			}

			cellValue = strings.TrimSpace(cellValue)
			if len(cellValue) == 0 {
				defaultValue := strings.TrimSpace(sheetData.DefaultValueList[j])
				cellValue = defaultValue
			}
			rowStream.WriteNodeValue(_cellType, cellValue)
		}

		fileStresam.WriteInt32(int32(i))
		fileStresam.WriteInt32(int32(rowStream.Len()))
		fileStresam.WriteRawBytes(rowStream.Buffer().Bytes())
	}
	fileStresam.WriteFile(src.GetBinFilePathName(sheetData.SheetName))
}

func init() {

}
