package common

import "strings"

type FieldType int

const (
	FieldType_Int32         FieldType = 1
	FieldType_Int32Array    FieldType = 2
	FieldType_Int32Array2D  FieldType = 3
	FieldType_UInt32        FieldType = 4
	FieldType_UInt32Array   FieldType = 5
	FieldType_UInt32Array2D FieldType = 6
	FieldType_Float         FieldType = 7
	FieldType_FloatArray    FieldType = 8
	FieldType_FloatArray2D  FieldType = 9
	FieldType_String        FieldType = 10
	FieldType_StringArray   FieldType = 11
	FieldType_StringArray2D FieldType = 12
	FieldType_Bool          FieldType = 13
	FieldType_BoolArray     FieldType = 14
	FieldType_BoolArray2D   FieldType = 15
)

func GetFieldType(strType string) FieldType {
	switch strType {
	case "int":
		return FieldType_Int32
		break
	case "int[]":
		return FieldType_Int32Array
		break
	case "int[][]":
		return FieldType_Int32Array2D
		break
	case "uint":
		return FieldType_UInt32
		break
	case "uint[]":
		return FieldType_UInt32Array
		break
	case "uint[][]":
		return FieldType_UInt32Array2D
		break
	case "float":
		return FieldType_Float
		break
	case "float[]":
		return FieldType_FloatArray
		break
	case "float[][]":
		return FieldType_FloatArray2D
		break
	case "bool":
		return FieldType_Bool
		break
	case "bool[]":
		return FieldType_BoolArray
		break
	case "bool[][]":
		return FieldType_BoolArray2D
		break
	case "string":
		return FieldType_String
		break
	case "string[]":
		return FieldType_StringArray
		break
	case "string[][]":
		return FieldType_StringArray2D
		break
	default:
		panic("unsupport type:" + strType)
	}
	return -1
}

func IsArray2D(fieldType string) bool {
	return strings.Contains(fieldType, "[][]")
}

func IsArray(fieldType string) bool {
	return strings.Contains(fieldType, "[]")
}

func GetBaseFieldTypeStr(fieldType string) string {
	return strings.ReplaceAll(fieldType, "[]", "")
}

func GetBaseFieldType(fieldType string) FieldType {
	return GetFieldType(GetBaseFieldTypeStr(fieldType))
}
