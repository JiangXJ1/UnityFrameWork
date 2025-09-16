package tool

import (
	"OneExcel/src"
	"strings"
)

// key 所在列
const l_key_column = 0

// 内容数据 所在列
const l_content_column = 2

// 备注 所在列
const l_des_column = 1

func GetLanguageData(sheetData src.SheetData) []src.LanguageData {
	var dataList []src.LanguageData

	for _, cellValues := range sheetData.AllRowValueList {
		//key
		key := strings.TrimSpace(cellValues[l_key_column])
		if len(key) == 0 {
			panic("KEY 不能为空！")
		}

		//内容
		content := strings.TrimSpace(cellValues[l_content_column])
		if len(content) == 0 {
			defaultValue := strings.TrimSpace(sheetData.DefaultValueList[l_content_column])
			content = defaultValue
		}

		//导出类型
		exportType := strings.TrimSpace(sheetData.FieldExportTypeList[l_content_column])
		if len(exportType) == 0 {
			defaultValue := strings.TrimSpace(sheetData.DefaultValueList[l_content_column])
			exportType = defaultValue
		}

		//描述
		des := cellValues[l_des_column]

		globalData := src.LanguageData{
			Key:        key,
			Content:    content,
			ExportType: exportType,
			Des:        des,
		}
		dataList = append(dataList, globalData)
	}

	return dataList
}
