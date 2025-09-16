package global

import (
	"strings"
)

var MaxFileCount int32 = 0
var CurFileCount int32 = 0
var CurTaskCount int32 = 0

const MIN_ROW_COUNT = 6
const LANGUAGE_NAME = "Language"

const EXCEL_DIR = "./Excels/"
const Export_DIR = "./Export/"
const CACHE_DIR = "./cache/"
const MD5_FILE = "./md5.txt"

const MIN_ROW = 5 //最少行数
const MIN_COL = 3 //最少列数
const GLOBAL_PREFIX = "Global_"
const LANGUAGE_PREFIX = "Language_"

const OUTPUT_DIR_BIN = "./Outputs/Client/CSharp/Bin/"
const OUTPUT_DIR_CS = "./Outputs/Client/CSharp/Cs/"
const OUTPUT_DIR_LUA = "./Outputs/Server//Lua/FromExcel/"
const OUTPUT_DIR_CSV = "./Outputs/Server/CSV/"

// 是否是全局表格
func IsGlobal(sheetName string) bool {
	return strings.HasPrefix(sheetName, GLOBAL_PREFIX) || strings.HasPrefix(sheetName, "@") || sheetName == "Global"
}

// 是否是语言表格
func IsLanguage(sheetName string) bool {
	return strings.HasPrefix(sheetName, LANGUAGE_PREFIX)
}
