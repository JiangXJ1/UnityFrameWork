package export_lua

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/tool"
	"OneExcel/src/fileUtil"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const Metatable = `

local #SCHEMA_NAME# = {}
local meta = {}
setmetatable(#SCHEMA_NAME#, meta)

function meta.Contains(t, key)
	return Data[key] ~= nil
end

function meta.__newindex(t, k,v)
	Log.Error("#SCHEMA_NAME# 表不允许修改 k=%s, v=%s", k, v)
end

function meta.__pairs(...)
	return pairs(Data)
end

function meta.__ipairs(...)
	return ipairs(Data)
end

function meta.__index(t,k)
	local val = Data[k]
	if not val then
		val = meta[k]
	end
	if not val then
		Log.Error("#SCHEMA_NAME# 表中未找到KEY:%s",k)
	end
	return val
end

function meta.__len(t)
	return #Data
end

`

type ClassLua struct {
}

func Write(sheetData src.SheetData) {
	needExportLua := NeedExportLua(sheetData)
	if !needExportLua {
		return
	}
	containServer := false
	for i := range sheetData.FieldExportTypeList {
		exportType := sheetData.FieldExportTypeList[i]
		if i > 0 {
			if !containServer && strings.Contains(exportType, "S") {
				containServer = true
			}
			if containServer {
				break
			}
		}
	}

	dir := ""
	filePath := ""
	if containServer {
		dir = src.OUTPUT_DIR_LUA_SERVER
		filePath = src.GetServerLuaFilePathName(sheetData.SheetName)
		WriteConfig(sheetData, dir, filePath, false)
	}
}

func WriteConfig(sheetData src.SheetData, dir string, filePath string, isClient bool) {
	fileUtil.CreateDir(dir)
	filePath = strings.ReplaceAll(filePath, "@", "")
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	isGlobal := sheetData.SheetName == "Schema_Global" || strings.HasPrefix(sheetData.SheetName, "Schema_@")
	WriteBegin(sheetData, write)
	if isGlobal {
		WriteGlobal(sheetData, write, isClient)
	} else if strings.Contains(sheetData.SheetName, "Language_") {
		WriteLanguage(sheetData, write, isClient)
	} else {
		WriteOneKey(sheetData, write, isClient)
	}

	WriteEnd(sheetData, write, isClient, isGlobal)
	write.Flush()
}

func WriteOneKey(sheetData src.SheetData, write *bufio.Writer, isClient bool) {
	write.WriteString(fmt.Sprintf("---@module %v\n", strings.ReplaceAll(sheetData.SheetName, "@", "")))
	for index, key := range sheetData.FieldNameList {
		if len(key) > 0 {
			exportType := sheetData.FieldExportTypeList[index]
			can_export := tool.CanExportLua(exportType, isClient)
			if !can_export {
				continue
			}
			strDes := ""
			if sheetData.DescList != nil && len(sheetData.DescList) > index {
				strDes = sheetData.DescList[index]
			}
			write.WriteString(fmt.Sprintf("---@field %v %v  @%v \n", key, sheetData.FieldTypeList[index], strings.ReplaceAll(strDes, "\n", " ")))
		}
	}
	//Write value
	if isClient {
		write.WriteString("\n---@type table<number,table>\n")
		write.WriteString("local Data =")
	} else {
		write.WriteString("local ")
		write.WriteString(strings.ReplaceAll(sheetData.SheetName, "@", ""))
		write.WriteString("=")
	}

	write.WriteString("\n{\n")

	rowLength := len(sheetData.AllRowValueList)

	exportTypeCount := len(sheetData.FieldExportTypeList)
	for i := range sheetData.AllRowValueList {
		cellValues := sheetData.AllRowValueList[i]
		count := len(cellValues)
		for j := range cellValues {
			if j >= exportTypeCount {
				continue
			}
			cellValue := cellValues[j]
			cellValue = strings.TrimSpace(cellValue)
			if j == 0 {
				if len(cellValue) == 0 {
					panic("KEY 不能为空！")
				}
				write.WriteString("\t[")
				write.WriteString(cellValue)
				write.WriteString("] = {")
			}
			exportType := sheetData.FieldExportTypeList[j]
			can_export := tool.CanExportLua(exportType, isClient)
			if !can_export {
				continue
			}
			_cellType := sheetData.FieldTypeList[j]
			_cellType = strings.TrimSpace(_cellType)
			if len(_cellType) == 0 {
				panic("字段类型不能为空字符串")
			}

			if len(cellValue) == 0 {
				defaultValue := strings.TrimSpace(sheetData.DefaultValueList[j])
				cellValue = defaultValue
			}

			key := sheetData.FieldNameList[j]
			write.WriteString(key)
			write.WriteString(" = ")

			fixCellValue, err := GetLuaCellValue(cellValue, _cellType)
			if err != nil {
				tip := sheetData.SheetName + " 行:" + strconv.Itoa(i) + " 列:" + strconv.Itoa(j)
				panic(err.Error() + tip)
			}
			write.WriteString(fixCellValue)

			if j != count-1 {
				write.WriteString(", ")
			}
		}
		if i != rowLength-1 {
			write.WriteString("},\n")
		} else {
			write.WriteString("}\n")
		}
	}
}

func WriteBegin(sheetData src.SheetData, write *bufio.Writer) {

}

func WriteEnd(sheetData src.SheetData, write *bufio.Writer, isClient bool, isGlobal bool) {
	write.WriteString("}\n")

	if isClient {
		str := strings.ReplaceAll(Metatable, "#SCHEMA_NAME#", strings.ReplaceAll(sheetData.SheetName, "@", ""))
		write.WriteString(str)
	}
	write.WriteString("\n")
	if !isGlobal {
		write.WriteString("---@type table<number,table>\n")
	}
	write.WriteString("return ")
	write.WriteString(strings.ReplaceAll(sheetData.SheetName, "@", ""))
}

func NeedExportLua(sheetData src.SheetData) bool {
	containName := false
	for i := range sheetData.FieldExportTypeList {
		if i > 0 {
			exportType := sheetData.FieldExportTypeList[i]
			if strings.Contains(exportType, "S") {
				containName = true
				break
			}
		}
	}
	return containName
}
