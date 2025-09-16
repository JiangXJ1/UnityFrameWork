package cs

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/common"
	"OneExcel/src/excelUtil/export_csharp"
	"OneExcel/src/fileUtil"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {

}

func GetTypeStr(strType string, name string) string {
	ft := common.GetFieldType(strType)
	switch ft {
	case common.FieldType_Bool:
		return "public readonly bool " + name + " = false;//*"
	case common.FieldType_BoolArray:
		return "public readonly ReadOnlyArray<bool> " + name + ";//*"
	case common.FieldType_BoolArray2D:
		return "public readonly ReadOnlyArray2D<bool> " + name + ";//*"
	case common.FieldType_Int32:
		return "public readonly int " + name + " = 0;"
	case common.FieldType_Int32Array:
		return "public readonly ReadOnlyArray<int> " + name + ";//*"
	case common.FieldType_Int32Array2D:
		return "public readonly ReadOnlyArray2D<int> " + name + ";//*"
	case common.FieldType_UInt32:
		return "public readonly uint " + name + " = 0;"
	case common.FieldType_UInt32Array:
		return "public readonly ReadOnlyArray<uint> " + name + ";//*"
	case common.FieldType_UInt32Array2D:
		return "public readonly ReadOnlyArray2D<uint> " + name + ";//*"
	case common.FieldType_Float:
		return "public readonly float " + name + " = 0;"
	case common.FieldType_FloatArray:
		return "public readonly ReadOnlyArray<float> " + name + ";//*"
	case common.FieldType_FloatArray2D:
		return "public readonly ReadOnlyArray2D<float> " + name + ";//*"
	case common.FieldType_String:
		return "public readonly string " + name + " = \"\";//*"
	case common.FieldType_StringArray:
		return "public readonly ReadOnlyArray<string> " + name + ";//*"
	case common.FieldType_StringArray2D:
		return "public readonly ReadOnlyArray2D<string> " + name + ";//*"
	default:
		src.IPanicEx("unsupport type:" + strType)
		panic("unsupport type:" + strType)
	}
}

func GetDataDefine(sheetData src.SheetData) string {
	str := ""
	cellValues := sheetData.AllRowValueList[0]
	count := len(cellValues)
	for j := range cellValues {
		//列
		key := sheetData.FieldNameList[j]
		_cellType := sheetData.FieldTypeList[j]
		content := GetTypeStr(_cellType, key)
		if j == count-1 {
			str = str + content
		} else {
			str = str + content + "\n\t\t"
		}
	}
	return str
}

func GetDeserializeArray(str string, name string, dateType string) string {
	temp := `uint count_#NAME# = reader.ReadUInt32();
			if(count_#NAME# > 0){
				#DATE_TYPE# [] tempArray = new #DATE_TYPE# [count_#NAME#];
				for (int i = 0; i < count_#NAME#; i++)
				{
					tempArray[i] = reader.#READ#;
				}
				#NAME# = new ReadOnlyArray<#DATE_TYPE#>(tempArray);
			}else
				#NAME# = default;`
	temp = strings.ReplaceAll(temp, "#READ#", str)
	temp = strings.ReplaceAll(temp, "#NAME#", name)
	temp = strings.ReplaceAll(temp, "#DATE_TYPE#", dateType)
	return temp
}
func GetDeserializeArray2D(str string, name string, dateType string) string {
	temp := `uint count_#NAME# = reader.ReadUInt32();
			if(count_#NAME# > 0){
				uint cellCount = 0;
				#DATE_TYPE# [][] tempArray = new #DATE_TYPE# [count_#NAME#][];
				for (int i = 0; i < count_#NAME#; i++)
				{
					cellCount = reader.ReadUInt32();
					tempArray[i] = new #DATE_TYPE# [cellCount];
					for (int j = 0; j < cellCount; j++)
					{
						tempArray[i][j] = reader.#READ#;
					}
				}
				#NAME# = new ReadOnlyArray2D<#DATE_TYPE#>(tempArray);
			}else
				#NAME# = default;`
	temp = strings.ReplaceAll(temp, "#READ#", str)
	temp = strings.ReplaceAll(temp, "#NAME#", name)
	temp = strings.ReplaceAll(temp, "#DATE_TYPE#", dateType)
	return temp
}
func GetDeserializeStr(strType string, name string) string {
	ft := common.GetFieldType(strType)
	switch ft {
	case common.FieldType_Bool:
		return name + " = reader.ReadUInt32() == 1;"
	case common.FieldType_BoolArray:
		return GetDeserializeArray("ReadUInt32() == 1", name, "bool")
	case common.FieldType_BoolArray2D:
		return GetDeserializeArray2D("ReadUInt32() == 1", name, "bool")
	case common.FieldType_Int32:
		return name + " = reader.ReadInt32();"
	case common.FieldType_Int32Array:
		return GetDeserializeArray("ReadInt32()", name, "int")
	case common.FieldType_Int32Array2D:
		return GetDeserializeArray2D("ReadInt32()", name, "int")
	case common.FieldType_UInt32:
		return name + " = reader.ReadUInt32();"
	case common.FieldType_UInt32Array:
		return GetDeserializeArray("ReadUInt32()", name, "uint")
	case common.FieldType_UInt32Array2D:
		return GetDeserializeArray2D("ReadUInt32()", name, "uint")
	case common.FieldType_Float:
		return name + " = reader.ReadFloat();"
	case common.FieldType_FloatArray:
		return GetDeserializeArray("ReadFloat()", name, "float")
	case common.FieldType_FloatArray2D:
		return GetDeserializeArray2D("ReadFloat()", name, "float")
	case common.FieldType_String:
		return name + " = reader.ReadString();"
	case common.FieldType_StringArray:
		return GetDeserializeArray("ReadString()", name, "string")
	case common.FieldType_StringArray2D:
		return GetDeserializeArray2D("ReadString()", name, "string")
	default:
		panic("unsupport type:" + strType)
	}
}

func GetDataDeserialize(sheetData src.SheetData) string {
	str := ""
	cellValues := sheetData.AllRowValueList[0]
	for j := range cellValues {
		//列
		key := sheetData.FieldNameList[j]
		_cellType := sheetData.FieldTypeList[j]
		content := GetDeserializeStr(_cellType, key)
		str = str + content + "\n\n\t\t\t"
	}
	return str
}

func WriteTableCS(sheetData src.SheetData) {
	containBin := export_csharp.NeedExportCSorBinary(sheetData)
	if containBin == false {
		fmt.Println(fmt.Sprintf("%v---%v:当前页签不包含c#配置", sheetData.FileName, sheetData.SheetName))
		return
	}

	if sheetData.SheetName == "Schema_Global" || strings.HasPrefix(sheetData.SheetName, "Schema_@") {
		WriteGlobalCS(sheetData)
	} else if strings.Contains(sheetData.SheetName, "Language_") {
		WriteLanguageCS(sheetData)
	} else {
		WriteCS(sheetData)
	}
	WriteEnum(sheetData)
}

func WriteCS(sheetData src.SheetData) {
	fileUtil.CreateDir(src.OUTPUT_DIR_CS)
	filePath := src.GetCSFilePathName(strings.ReplaceAll(sheetData.SheetName, "@", ""))
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	str_define := ""
	str_deserialize := ""
	count := len(sheetData.FieldExportTypeList)
	for j := range sheetData.FieldExportTypeList {
		key := sheetData.FieldNameList[j]
		_cellType := sheetData.FieldTypeList[j]

		exportType := sheetData.FieldExportTypeList[j]
		if !strings.Contains(exportType, "C") {
			continue
		}

		content_define := GetTypeStr(_cellType, key)
		content_deser := GetDeserializeStr(_cellType, key)
		if j == count-1 {
			str_define = str_define + content_define
			str_deserialize = str_deserialize + content_deser
		} else {
			str_define = str_define + content_define + "\n\t\t"
			str_deserialize = str_deserialize + content_deser + "\n\n\t\t\t"
		}
	}

	str_dataDefine := str_define
	str_dataDeserialize := str_deserialize

	content := CsharpTemplate

	content = strings.ReplaceAll(content, "#KEY_TYPE#", sheetData.FieldTypeList[0])
	content = strings.ReplaceAll(content, "#KEY_NAME#", sheetData.FieldNameList[0])
	content = strings.ReplaceAll(content, "#DATA_DEFINE#", str_dataDefine)
	content = strings.ReplaceAll(content, "#DESERIALIZE#", str_dataDeserialize)
	content = strings.ReplaceAll(content, "#CLASS_NAME#", strings.ReplaceAll(sheetData.SheetName, "@", ""))

	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
}
