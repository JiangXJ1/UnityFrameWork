package src

import (
	"bytes"
	"os"
)

const EXCEL_ONE_PATH = "./Excels/UITextDataAuto.xlsx"
const EXCEL_MD5_PATH = "./MD5.txt"
const EXCEL_DIR = "./Excels/"
const HEAD_SHEET_INGORE = "#"
const HEAD_Row_INGORE = "#"
const HEAD_CELL_KEY = "*"

// 描述行
const DES_ROW = 4

type SheetData struct {
	FileName            string
	SheetName           string
	DescList            []string
	FieldTypeList       []string
	FieldNameList       []string
	DefaultValueList    []string
	FieldExportTypeList []string

	AllRowValueList [][]string
}

type GlobalData struct {
	Key         string
	ContentType string
	ExportType  string
	Content     string
	Des         string
}

type LanguageData struct {
	Key        string
	ExportType string
	Content    string
	Des        string
}

func GetExcelFilePathName(fileNameWithExt string) string {
	var buffer bytes.Buffer
	buffer.WriteString(EXCEL_DIR)
	buffer.WriteString(fileNameWithExt)
	return buffer.String()
}

// Lua
const OUTPUT_DIR_LUA_SERVER = "./Outputs/Server/Lua/FromExcel/"
const OUTPUT_DIR_BIN = "./Outputs/Client/CSharp/Bin/"
const OUTPUT_DIR_CS = "./Outputs/Client/CSharp/Cs/"
const OUTPUT_DIR_CSV = "./Outputs/Server/CSV/"

func ClearPath() {
	os.RemoveAll(OUTPUT_DIR_LUA_SERVER)
	os.RemoveAll(OUTPUT_DIR_BIN)
	os.RemoveAll(OUTPUT_DIR_CS)
}

func GetServerLuaFilePathName(fileName string) string {
	var buffer bytes.Buffer
	buffer.WriteString(OUTPUT_DIR_LUA_SERVER)
	buffer.WriteString(fileName)
	buffer.WriteString(".lua")
	return buffer.String()
}

func GetBinFilePathName(fileName string) string {
	var buffer bytes.Buffer
	buffer.WriteString(OUTPUT_DIR_BIN)
	buffer.WriteString(fileName)
	buffer.WriteString(".bytes")
	return buffer.String()
}

func GetCSFilePathName(fileName string) string {
	var buffer bytes.Buffer
	buffer.WriteString(OUTPUT_DIR_CS)
	buffer.WriteString(fileName)
	buffer.WriteString(".cs")
	return buffer.String()
}

func GetCSVFilePathName(fileName string) string {
	var buffer bytes.Buffer
	buffer.WriteString(OUTPUT_DIR_CSV)
	buffer.WriteString(fileName)
	buffer.WriteString(".csv")
	return buffer.String()
}
