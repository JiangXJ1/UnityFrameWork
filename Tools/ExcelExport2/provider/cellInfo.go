package provider

import (
	"demo-go-excel/util/common"
)

type CellInfo struct {
	FieldInfo FieldInfo
	Content   string
	Value     interface{}
}

func (c *CellInfo) TryParse() {
	if !c.FieldInfo.Valid {
		return
	}
	switch c.FieldInfo.FieldType {
	case "string", "string[]", "string[][]":
		c.Value = common.ParseString(c.FieldInfo.FieldType, c.Content)
	case "bool", "bool[]", "bool[][]":
		c.Value = common.ParseBool(c.FieldInfo.FieldType, c.Content)
	case "int", "int[]", "int[][]":
		c.Value = common.ParseInt32(c.FieldInfo.FieldType, c.Content)
	case "uint", "uint[]", "uint[][]":
		c.Value = common.ParseUInt32(c.FieldInfo.FieldType, c.Content)
	case "float", "float[]", "float[][]":
		c.Value = common.ParseFloat(c.FieldInfo.FieldType, c.Content)
	default:
		panic("未知的类型: " + c.FieldInfo.FieldType)
	}
}
