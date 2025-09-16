package tool

import (
	"OneExcel/src"
	"strings"
)

// key 所在列
const key_column = 0

// 内容数据类型 所在列
const contentType_column = 1

// 内容数据导出规则 所在列
const contentExportType_column = 2

// 内容数据 所在列
const content_column = 3

// 备注 所在列
const des_column = 4

func GetGlobalContent(sheetData src.SheetData) []src.GlobalData {
	var dataList []src.GlobalData

	for _, cellValues := range sheetData.AllRowValueList {
		//key
		key := strings.TrimSpace(cellValues[key_column])
		if len(key) == 0 {
			panic("KEY 不能为空！")
		}

		//数据类型
		contentType := strings.TrimSpace(cellValues[contentType_column])
		if len(contentType) == 0 {
			defaultValue := strings.TrimSpace(sheetData.DefaultValueList[contentType_column])
			contentType = defaultValue
		}

		//导出类型
		exportType := strings.TrimSpace(cellValues[contentExportType_column])
		if len(exportType) == 0 {
			defaultValue := strings.TrimSpace(sheetData.DefaultValueList[contentExportType_column])
			exportType = defaultValue
		}

		//内容
		content := strings.TrimSpace(cellValues[content_column])
		if len(content) == 0 {
			defaultValue := strings.TrimSpace(sheetData.DefaultValueList[content_column])
			content = defaultValue
		}

		//描述
		des := cellValues[des_column]

		globalData := src.GlobalData{
			Key:         key,
			ContentType: contentType,
			ExportType:  exportType,
			Content:     content,
			Des:         des,
		}
		dataList = append(dataList, globalData)
	}

	return dataList
}
