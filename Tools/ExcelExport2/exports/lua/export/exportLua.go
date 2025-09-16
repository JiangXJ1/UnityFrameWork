package export

import (
	"bytes"
	"demo-go-excel/util/common"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func GetLuaCellValue(cellValue string, cellType string) (string, error) {
	ft := common.GetFieldType(cellType)
	var buffer bytes.Buffer
	switch ft {
	case common.FieldType_Int32, common.FieldType_Float, common.FieldType_UInt32:
		return cellValue, nil
	case common.FieldType_Int32Array, common.FieldType_FloatArray, common.FieldType_UInt32Array:
		return GetNumberArray(cellValue, buffer), nil
	case common.FieldType_Int32Array2D, common.FieldType_FloatArray2D, common.FieldType_UInt32Array2D:
		return GetNumberArray2D(cellValue, buffer), nil
	case common.FieldType_Bool:
		return GetBool(cellValue), nil
	case common.FieldType_BoolArray:
		return GetBoolArray(cellValue, buffer), nil
	case common.FieldType_BoolArray2D:
		return GetBoolArray2D(cellValue, buffer), nil
	case common.FieldType_String:
		return GetString(cellValue, buffer), nil
	case common.FieldType_StringArray:
		return GetStringArray(cellValue, buffer), nil
	case common.FieldType_StringArray2D:
		return GetStringArray2D(cellValue, buffer), nil
	default:
		log.Fatalf("Unknow type: %s ", cellType)
		panic(fmt.Sprintf("Unknow type: %s ", cellType))
		return cellValue, nil
	}
	return cellValue, nil
}

func GetBool(cellValue string) string {
	value, err := strconv.Atoi(cellValue)
	if err != nil || (value != 0 && value != 1) {
		fmt.Printf("bool值配置不正确！")
	}
	if value == 1 {
		return "true"
	} else {
		return "false"
	}
}

func GetBoolArray(cellValue string, buffer bytes.Buffer) string {
	if len(cellValue) == 0 {
		buffer.WriteString("{}")
	} else {
		cellArray := strings.Split(cellValue, "*")
		buffer.WriteString("{")
		length := len(cellArray)
		length--
		for i := range cellArray {
			num := cellArray[i]
			buffer.WriteString(GetBool(num))
			if i != length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
	}
	return buffer.String()
}

func GetBoolArray2D(cellValue string, buffer bytes.Buffer) string {
	if len(cellValue) == 0 {
		buffer.WriteString("{}")
	} else {
		cellArray := strings.Split(cellValue, "|")
		buffer.WriteString("{")
		length := len(cellArray)
		length--
		var buffer2 bytes.Buffer
		for i := range cellArray {
			buffer2.Reset()
			content := GetBoolArray(cellArray[i], buffer2)
			buffer.WriteString(content)
			if i != length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
	}
	return buffer.String()
}

func GetString(cellValue string, buffer bytes.Buffer) string {
	buffer.WriteString("'")
	buffer.WriteString(cellValue)
	buffer.WriteString("'")
	return buffer.String()
}

func GetNumberArray(cellValue string, buffer bytes.Buffer) string {
	if len(cellValue) == 0 {
		buffer.WriteString("{}")
	} else {
		cellArray := strings.Split(cellValue, "*")
		buffer.WriteString("{")
		length := len(cellArray)
		length--
		for i := range cellArray {
			num := cellArray[i]
			buffer.WriteString(num)
			if i != length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
	}
	return buffer.String()
}
func GetNumberArray2D(cellValue string, buffer bytes.Buffer) string {
	//1:2:3|4:5
	//{{1,2,3},{4,5}}
	if len(cellValue) == 0 {
		buffer.WriteString("{}")
	} else {
		cellArray := strings.Split(cellValue, "|")
		buffer.WriteString("{")
		length := len(cellArray)
		length--
		var buffer2 bytes.Buffer
		for i := range cellArray {
			buffer2.Reset()
			content := GetNumberArray(cellArray[i], buffer2)
			buffer.WriteString(content)
			if i != length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
	}
	return buffer.String()
}

func GetStringArray(cellValue string, buffer bytes.Buffer) string {
	//"a":"b":"c":"
	//{"a","b","c"}
	if len(cellValue) == 0 {
		buffer.WriteString("{}")
	} else {
		cellArray := strings.Split(cellValue, "*")
		buffer.WriteString("{")
		length := len(cellArray)
		length--
		var buffer2 bytes.Buffer
		for i := range cellArray {
			num := cellArray[i]
			buffer2.Reset()
			buffer.WriteString(GetString(num, buffer2))
			if i != length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
	}
	return buffer.String()
}
func GetStringArray2D(cellValue string, buffer bytes.Buffer) string {
	//a:b|c
	//{{"a","b"},{"c"}}

	if len(cellValue) == 0 {
		buffer.WriteString("{}")
	} else {
		cellArray := strings.Split(cellValue, "|")
		buffer.WriteString("{")
		length := len(cellArray)
		length--
		var buffer2 bytes.Buffer
		for i := range cellArray {
			buffer2.Reset()
			content := GetStringArray(cellArray[i], buffer2)
			buffer.WriteString(content)
			if i != length {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
	}
	return buffer.String()
}
