package export

import (
	"bytes"
	"demo-go-excel/util/common"
	"fmt"
	"text/template"
)

type CodeField struct {
	ReadMethod      string
	Name            string
	DataType        string
	ReadArrayMethod string
}

type ReadCode struct {
	ReadMethod      string
	DataType        string
	IsArray         bool
	IsArray2D       bool
	ReadArrayMethod string
}

func init() {

}

func GetTypeStr(strType string, name string, comment string) string {
	fieldTemplates := map[common.FieldType]string{
		common.FieldType_Bool:          "public readonly bool %s = false;//%s",
		common.FieldType_BoolArray:     "public readonly ReadOnlyArray<bool> %s;//%s",
		common.FieldType_BoolArray2D:   "public readonly ReadOnlyArray2D<bool> %s;//%s",
		common.FieldType_Int32:         "public readonly int %s = 0;//%s",
		common.FieldType_Int32Array:    "public readonly ReadOnlyArray<int> %s;//%s",
		common.FieldType_Int32Array2D:  "public readonly ReadOnlyArray2D<int> %s;//%s",
		common.FieldType_UInt32:        "public readonly uint %s = 0;//%s",
		common.FieldType_UInt32Array:   "public readonly ReadOnlyArray<uint> %s;//%s",
		common.FieldType_UInt32Array2D: "public readonly ReadOnlyArray2D<uint> %s;//%s",
		common.FieldType_Float:         "public readonly float %s = 0;//%s",
		common.FieldType_FloatArray:    "public readonly ReadOnlyArray<float> %s;//%s",
		common.FieldType_FloatArray2D:  "public readonly ReadOnlyArray2D<float> %s;//%s",
		common.FieldType_String:        "public readonly string %s = \"\";//%s",
		common.FieldType_StringArray:   "public readonly ReadOnlyArray<string> %s;//%s",
		common.FieldType_StringArray2D: "public readonly ReadOnlyArray2D<string> %s;//%s",
	}

	ft := common.GetFieldType(strType)
	if template, ok := fieldTemplates[ft]; ok {
		return fmt.Sprintf(template, name, comment)
	}

	panic("unsupported type: " + strType)
}

func GetDeserializeArray(str string, name string, dataType string, readArrayMethod string) string {
	const tmpl = `uint count_{{.Name}} = reader.ReadUInt32();
			{{.Name}} = reader.{{.ReadArrayMethod}}(count_{{.Name}});`
	params := CodeField{
		ReadMethod:      str,
		Name:            name,
		DataType:        dataType,
		ReadArrayMethod: readArrayMethod,
	}
	var buf bytes.Buffer
	t := template.Must(template.New("deserialize").Parse(tmpl))
	if err := t.Execute(&buf, params); err != nil {
		panic("template execution failed: " + err.Error())
	}
	return buf.String()
}
func GetDeserializeArray2D(str string, name string, dataType string, readArrayMethod string) string {
	const tmpl = `uint count_{{.Name}} = reader.ReadUInt32();
			{{.Name}} = reader.{{.ReadArrayMethod}}(count_{{.Name}});`

	// 定义模板参数
	params := CodeField{
		ReadMethod:      str,
		Name:            name,
		DataType:        dataType,
		ReadArrayMethod: readArrayMethod,
	}

	// 渲染模板
	var buf bytes.Buffer
	t := template.Must(template.New("deserialize2D").Parse(tmpl))
	if err := t.Execute(&buf, params); err != nil {
		panic("template execution failed: " + err.Error())
	}
	return buf.String()
}
func GetDeserializeStr(strType string, name string) string {
	fieldHandlers := map[common.FieldType]ReadCode{
		common.FieldType_Bool:          {"ReadUInt32() == 1", "bool", false, false, ""},
		common.FieldType_BoolArray:     {"ReadUInt32() == 1", "bool", true, false, "ReadBoolArray"},
		common.FieldType_BoolArray2D:   {"ReadUInt32() == 1", "bool", false, true, "ReadBoolArray2D"},
		common.FieldType_Int32:         {"ReadInt32()", "int", false, false, ""},
		common.FieldType_Int32Array:    {"ReadInt32()", "int", true, false, "ReadInt32Array"},
		common.FieldType_Int32Array2D:  {"ReadInt32()", "int", false, true, "ReadInt32Array2D"},
		common.FieldType_UInt32:        {"ReadUInt32()", "uint", false, false, ""},
		common.FieldType_UInt32Array:   {"ReadUInt32()", "uint", true, false, "ReadUInt32Array"},
		common.FieldType_UInt32Array2D: {"ReadUInt32()", "uint", false, true, "ReadUInt32Array2D"},
		common.FieldType_Float:         {"ReadFloat()", "float", false, false, ""},
		common.FieldType_FloatArray:    {"ReadFloat()", "float", true, false, "ReadFloatArray"},
		common.FieldType_FloatArray2D:  {"ReadFloat()", "float", false, true, "ReadFloatArray2D"},
		common.FieldType_String:        {"ReadString()", "string", false, false, ""},
		common.FieldType_StringArray:   {"ReadString()", "string", true, false, "ReadStringArray"},
		common.FieldType_StringArray2D: {"ReadString()", "string", false, true, "ReadStringArray2D"},
	}

	ft := common.GetFieldType(strType)
	handler, ok := fieldHandlers[ft]
	if !ok {
		panic("unsupported type: " + strType)
	}

	// 根据类型生成对应代码
	if handler.IsArray2D {
		return GetDeserializeArray2D(handler.ReadMethod, name, handler.DataType, handler.ReadArrayMethod)
	} else if handler.IsArray {
		return GetDeserializeArray(handler.ReadMethod, name, handler.DataType, handler.ReadArrayMethod)
	} else {
		return name + " = reader." + handler.ReadMethod + ";"
	}
}
