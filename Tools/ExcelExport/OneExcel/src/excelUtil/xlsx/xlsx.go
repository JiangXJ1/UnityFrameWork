package xlsx

import (
	"OneExcel/src"
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
	"github.com/wxnacy/wgo/arrays"
)

func GetData(filepath string, fileName string) []src.SheetData {
	//timeNow := time.Now().UnixNano()
	var results []src.SheetData

	src.CurrentFileName = fileName
	f, err := xlsx.OpenFile(filepath)
	if err != nil {
		fmt.Println("打开异常：" + fileName)
		fmt.Println(err)
		return results
	}

	// Get sheet
	for _, sheet := range f.Sheets {
		src.CurrentSheetName = sheet.Name
		sheetName := sheet.Name
		if strings.HasPrefix(sheetName, src.HEAD_SHEET_INGORE) {
			continue
		}
		var descList []string
		var fieldTypeList []string
		var fieldNameList []string
		var fieldExportTypeList []string
		var defaultList []string
		var allRowValueList [][]string
		//var rowIndexIgnore []int
		var cellIndexIgnore []int

		var tempKes []string
		var realIndex int
		columnCount := len(sheet.Rows[0].Cells)
		//firstRow := sheet.Rows[0].Cells
		//for i := len(firstRow) - 1; i >= 0; i-- {
		//	if len(firstRow[i].String())>0 {
		//		columnCount = i+1
		//		break
		//	}
		//}
		for rowIndex, row := range sheet.Rows {
			src.CurrentRow = rowIndex
			src.CurrentColum = 1
			sheet.Row(rowIndex)
			var cellValueList []string
			if len(row.Cells) <= 1 {
				continue
			}
			temp := row.Cells[0].String()
			key := row.Cells[1].String()
			if temp == src.HEAD_Row_INGORE || (rowIndex > 3 && len(key) == 0) {
				continue
			}
			if rowIndex > 4 && arrays.Contains(tempKes, key) >= 0 {
				src.PrintTableInfo("Key重复：" + key)
				panic("")
			}
			tempKes = append(tempKes, key)

			for cellIndex := 0; cellIndex < columnCount; cellIndex++ {
				if cellIndex == 0 {
					continue
				}
				realIndex = cellIndex - 1
				src.CurrentColum = cellIndex
				cellValue := ""
				if cellIndex < len(row.Cells) {
					cellValue = row.Cells[cellIndex].String()
				}

				if arrays.Contains(cellIndexIgnore, cellIndex) >= 0 {
					continue
				}
				if rowIndex == 0 { //导出规则：客户端、服务器、c#、lua
					if len(fieldExportTypeList) == 0 { //KEY
						fieldExportTypeList = append(fieldExportTypeList, "C/S")
					} else {
						fieldExportTypeList = append(fieldExportTypeList, cellValue)
					}

				} else if rowIndex == 1 { //字段类型
					if fieldExportTypeList != nil && len(fieldExportTypeList) > len(fieldTypeList) {
						fieldTypeList = append(fieldTypeList, cellValue)
						if len(cellValue) == 0 && len(fieldExportTypeList[len(fieldTypeList)-1]) > 0 {
							src.IPanic("当前列字段类型不能为空！")
						}
					}
				} else if rowIndex == 2 { //字段名称
					if len(fieldExportTypeList[realIndex]) > 0 && len(cellValue) > 0 && arrays.Contains(fieldNameList, cellValue) >= 0 {
						src.IPanic("字段名重复：" + cellValue)
					}
					fieldNameList = append(fieldNameList, cellValue)
				} else if rowIndex == 3 { //默认值
					defaultList = append(defaultList, cellValue)
				} else if rowIndex == src.DES_ROW { //描述
					descList = append(descList, cellValue)
				} else {
					if len(fieldExportTypeList) > cellIndex && len(fieldExportTypeList[realIndex]) > 0 {
						src.CurrentFiledName = sheet.Rows[2].Cells[cellIndex].String()
						//if len(cellValue) > 0 {
						CheckLegitimacy(cellValue, fieldTypeList[realIndex])
						//}else {
						//	CheckLegitimacy(defaultList[realIndex], fieldTypeList[realIndex])
						//}
					}
					cellValueList = append(cellValueList, cellValue)
				}
			}

			if len(cellValueList) <= 0 {
				continue
			}
			allRowValueList = append(allRowValueList, cellValueList)
		}
		if fieldExportTypeList == nil {
			fmt.Println("s")
		}
		sheetData := src.SheetData{
			FileName:            fileName,
			SheetName:           fmt.Sprintf("Schema_%v", sheetName),
			DescList:            descList,
			FieldTypeList:       fieldTypeList,
			FieldNameList:       fieldNameList,
			AllRowValueList:     allRowValueList,
			FieldExportTypeList: fieldExportTypeList,
			DefaultValueList:    defaultList,
		}
		var dataList = CheckSplit(sheetData)
		for _, value := range dataList {
			results = append(results, value)
		}
	}
	//timeNow = time.Now().UnixNano() - timeNow
	//timeNow = timeNow / 1e6
	//fmt.Printf("\nTotal time : %d ms. \n", timeNow)
	return results
}

// 检查分表(表名包含 "_split")
func CheckSplit(sheetData src.SheetData) []src.SheetData {
	var data []src.SheetData
	if !strings.Contains(sheetData.SheetName, "_split") {
		fmt.Println(sheetData.SheetName)
		data = append(data, sheetData)
		return data
	} else {
		var temp []string

		for index, cellName := range sheetData.FieldNameList {
			if strings.Contains(cellName, "_*") && len(sheetData.FieldExportTypeList[index]) > 0 {
				strs := strings.Split(cellName, "_*")
				if len(strs) > 1 {
					id := arrays.ContainsString(temp, strs[1])
					if id == -1 {
						temp = append(temp, strs[1])
					}
				}
			}
		}
		if len(temp) == 0 {
			data = append(data, sheetData)
			return data
		} else {
			for _, name := range temp {
				var descList []string
				var fieldTypeList []string
				var fieldNameList []string
				var filedExportTypeList []string
				var defaultList []string
				var allRowValueList [][]string

				for index, fieldName := range sheetData.FieldNameList {
					if !strings.Contains(fieldName, "_*") || strings.Contains(fieldName, "_*"+name) {
						descList = append(descList, sheetData.DescList[index])
						fieldTypeList = append(fieldTypeList, sheetData.FieldTypeList[index])
						if strings.Contains(fieldName, "_*") {
							strs := strings.Split(fieldName, "_*")
							fieldNameList = append(fieldNameList, strs[0])
						} else {
							fieldNameList = append(fieldNameList, fieldName)
						}
						filedExportTypeList = append(filedExportTypeList, sheetData.FieldExportTypeList[index])
						defaultList = append(defaultList, sheetData.DefaultValueList[index])
					}
				}

				for _, value := range sheetData.AllRowValueList {
					var cellValueList []string
					for cellIndex, cell := range value {
						if cellIndex < len(sheetData.FieldExportTypeList) {
							fieldName := sheetData.FieldNameList[cellIndex]
							if !strings.Contains(fieldName, "_*") || strings.Contains(fieldName, "_*"+name) {
								cellValueList = append(cellValueList, cell)
							}
						}
					}
					allRowValueList = append(allRowValueList, cellValueList)
				}

				tempSheetData := src.SheetData{
					FileName:            sheetData.FileName,
					SheetName:           strings.ReplaceAll(sheetData.SheetName, "split", name),
					DescList:            descList,
					FieldTypeList:       fieldTypeList,
					FieldNameList:       fieldNameList,
					AllRowValueList:     allRowValueList,
					FieldExportTypeList: filedExportTypeList,
					DefaultValueList:    defaultList,
				}

				data = append(data, tempSheetData)
			}
		}
		return data
	}
}

func CheckFieldType() {

}
