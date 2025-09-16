package xlsx

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/common"
	"strconv"
	"strings"
)

func SplitArray(str string) []string {
	strArray := strings.Split(str, "*")
	return strArray
}

func GetArray(value string) []string {
	array := SplitArray(value)
	return array
}

func GetArray2D(value string) []string {
	strArray := strings.Split(value, "|")
	return strArray
}

func checkValue(fieldType string, value string) {
	ft := common.GetFieldType(fieldType)
	switch ft {
	case common.FieldType_Int32:
		CheckInt32(value)
		break
	case common.FieldType_UInt32:
		CheckUInt32(value)
		break
	case common.FieldType_Float:
		CheckFloat(value)
		break
	case common.FieldType_Bool:
		CheckBool(value)
		break
	case common.FieldType_String:
		CheckString(value)
		break
	default:
		src.IPanic("当前字段导出类型不支持！" + fieldType)
	}
}

// 检查内容合法行
func CheckLegitimacy(value string, fieldType string) {
	//配置内容和字段定义类型不符
	if common.IsArray2D(fieldType) {
		if len(value) > 0 {
			arrayRoot := GetArray2D(value)
			for i := range arrayRoot {
				array := GetArray(arrayRoot[i])
				for _, cell := range array {
					checkValue(common.GetBaseFieldTypeStr(fieldType), cell)
				}
			}
		}
	} else if common.IsArray(fieldType) {

		if len(value) > 0 {
			array := GetArray(value)
			for _, cell := range array {
				checkValue(common.GetBaseFieldTypeStr(fieldType), cell)
			}
		}
	} else {
		checkValue(fieldType, value)
	}
}

func CheckBool(value string) bool {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil || (v != 0 && v != 1) {
		src.IPanic("可能配错!bool 必须填1或者0 " + value)
	}
	return true
}

func CheckInt32(value string) bool {
	_, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		src.IPanic("可能配错!无法转换成int: " + value)
	}
	return true
}
func CheckUInt32(value string) bool {
	_, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		src.IPanic("可能配错!无法转换成uint: " + value)
	}
	return true
}
func CheckFloat(value string) bool {
	_, err := strconv.ParseFloat(value, 32)
	if err != nil {
		src.IPanic("可能配错!无法转换成float: " + value)
	}
	return true
}
func CheckString(value string) bool {
	return true
}
