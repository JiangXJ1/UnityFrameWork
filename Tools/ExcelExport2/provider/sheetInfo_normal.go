package provider

import (
	"bufio"
	"demo-go-excel/exports/cs/export"
	exportLua "demo-go-excel/exports/lua/export"
	"demo-go-excel/util"
	"demo-go-excel/util/common"
	"demo-go-excel/util/fileUtil"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type SheetInfoNormal struct {
	SheetInfoBase
}

func (s *SheetInfoNormal) Read(task *Task) {
	s.FileName = task.FileName
	s.Sheet = task.Sheet
	//解析字段信息
	s.GetFieldInfo()
	//解析内容
	s.GetContent()
}

func GetCell() {

}

// 解析字段信息
func (s *SheetInfoNormal) GetFieldInfo() {
	columnCount := len(s.Sheet.Rows[0].Cells)
	s.FieldInfos = make([]FieldInfo, columnCount)
	names := make(map[string]bool)
	for i := 0; i < columnCount; i++ {
		exportType := s.Sheet.Rows[0].Cells[i].String()

		fieldType := ""
		fieldName := ""
		defaultValue := ""
		comment := ""
		if len(s.Sheet.Rows[1].Cells) > i {
			fieldType = s.Sheet.Rows[1].Cells[i].String()
		}
		if len(s.Sheet.Rows[2].Cells) > i {
			fieldName = s.Sheet.Rows[2].Cells[i].String()
		}
		if len(s.Sheet.Rows[3].Cells) > i {
			defaultValue = s.Sheet.Rows[3].Cells[i].String()
		}
		if len(s.Sheet.Rows[4].Cells) > i {
			comment = s.Sheet.Rows[4].Cells[i].String()
		}

		if len(exportType) > 0 && len(fieldName) > 0 && names[fieldName] {
			fmt.Println("字段名重复：" + fieldName)
			panic("字段名重复：" + fieldName)
		}
		names[fieldName] = true
		fieldInfo := FieldInfo{
			Index:        i,
			Valid:        i > 0 && len(exportType) > 0 && len(fieldType) > 0 && len(fieldName) > 0,
			ExportType:   strings.TrimSpace(exportType),
			Name:         strings.TrimSpace(fieldName),
			FieldType:    strings.TrimSpace(fieldType),
			DefaultValue: strings.TrimSpace(defaultValue),
			Comment:      strings.TrimSpace(strings.ReplaceAll(comment, "\n", " ")),
		}
		s.FieldInfos[i] = fieldInfo
	}
	for i, item := range s.FieldInfos {
		if i < 2 || !item.Valid {
			continue
		}
		if !s.CanExportC && item.SupportClient() {
			s.CanExportC = true
		}
		if !s.CanExportS && item.SupportServer() {
			s.CanExportS = true
		}
		if s.CanExportC && s.CanExportS {
			break
		}
	}
}

// 解析内容
func (s *SheetInfoNormal) GetContent() {
	columnCount := len(s.Sheet.Rows[0].Cells)
	s.CellInfos = make([][]CellInfo, len(s.Sheet.Rows)-5)
	for rowIndex := 5; rowIndex < len(s.Sheet.Rows); rowIndex++ {
		row := s.Sheet.Rows[rowIndex]
		if len(row.Cells) < 2 || len(strings.TrimSpace(row.Cells[1].Value)) == 0 || strings.HasPrefix(strings.TrimSpace(row.Cells[0].Value), "#") {
			continue
		}
		rowData := make([]CellInfo, columnCount)
		for colIndex := 0; colIndex < columnCount; colIndex++ {
			if !s.FieldInfos[colIndex].Valid {
				rowData[colIndex] = CellInfo{}
				continue
			}

			var str string
			if colIndex >= len(row.Cells) {
				str = ""
			} else {
				str = strings.TrimSpace(row.Cells[colIndex].String())
			}

			if len(str) == 0 && len(s.FieldInfos[colIndex].DefaultValue) > 0 {
				str = s.FieldInfos[colIndex].DefaultValue
			}

			cellInfo := CellInfo{
				FieldInfo: s.FieldInfos[colIndex],
				Content:   str,
			}

			// 添加错误处理
			func() {
				defer func() {
					if r := recover(); r != nil {
						panic(NewExportError(
							s.FileName,
							s.Sheet.Name,
							s.FieldInfos[colIndex].Name,
							rowIndex+1,
							fmt.Sprint(r),
						))
					}
				}()
				cellInfo.TryParse()
			}()

			rowData[colIndex] = cellInfo
			s.CellInfos[rowIndex-5] = rowData
		}
	}
}

// Export 导出
func (s *SheetInfoNormal) Export() {
	if s.CanExportC {
		s.ExportCsharp()
		s.ExportBin()
	}
	if s.CanExportS {
		s.ExportLua()
		s.ExportCSV()
	}
}

// ExportCsharp 导出C#代码
func (s *SheetInfoNormal) ExportCsharp() {
	filePath := fileUtil.GetCSFilePathName(s.Sheet.Name)
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(file)

	strDefine := ""
	strDeserialize := ""
	count := len(s.FieldInfos)
	for j := range s.FieldInfos {
		fieldInfo := s.FieldInfos[j]
		if !fieldInfo.SupportCSharp() {
			continue
		}
		key := fieldInfo.Name
		fieldType := fieldInfo.FieldType
		contentDefine := export.GetTypeStr(fieldType, key, fieldInfo.Comment)
		contentDeser := export.GetDeserializeStr(fieldType, key)
		if j == count-1 {
			strDefine = strDefine + contentDefine
			strDeserialize = strDeserialize + contentDeser
		} else {
			strDefine = strDefine + contentDefine + "\n\t\t"
			strDeserialize = strDeserialize + contentDeser + "\n\n\t\t\t"
		}
	}
	//模版参数
	normalTemplate := export.CsharpNormalTemplateParam{
		ClassName:   fmt.Sprintf("Schema_%s", strings.ReplaceAll(s.Sheet.Name, "@", "")),
		KeyName:     s.FieldInfos[1].Name,
		KeyType:     s.FieldInfos[1].FieldType,
		DataDefine:  strDefine,
		Deserialize: strDeserialize,
	}
	//生成代码
	content := normalTemplate.GenerateCsharpTemplate()
	//写入文件
	write := bufio.NewWriter(file)
	util.WriteStringToFile(write, content)
	err = write.Flush()
	if err != nil {
		return
	}
}

// ExportBin 导出二进制文件
func (s *SheetInfoNormal) ExportBin() {
	fileStream := common.NewStream()
	count := len(s.FieldInfos)
	index := 0
	for _, cellValues := range s.CellInfos {
		if cellValues == nil || len(cellValues) < 3 {
			continue
		}
		rowStream := common.NewStream()
		//行
		for j, cellValue := range cellValues {
			if j == 0 || j >= count || j > 1 && !s.FieldInfos[j].SupportCSharp() {
				continue
			}
			//列
			fieldType := cellValue.FieldInfo.FieldType

			rowStream.WriteNodeValue(fieldType, cellValue.Value)
			//fmt.Printf("------------------：%v -- %v\n", j, cellValue.Value)
		}

		fileStream.WriteInt32(int32(index))
		fileStream.WriteInt32(int32(rowStream.Len()))
		fileStream.WriteRawBytes(rowStream.Buffer().Bytes())
		//fmt.Printf(">>>>>>>>>：%v -- %v -- %v\n", index, int32(rowStream.Len()), len(rowStream.Buffer().Bytes()))
		index++
	}

	err := fileStream.WriteFile(fileUtil.GetBinFilePathName(s.Sheet.Name))
	if err != nil {
		fmt.Printf("导出二进制文件失败：%s\n", err.Error())
	}
}

func (s *SheetInfoNormal) ExportLua() {
	sheetName := s.Sheet.Name
	filePath := fileUtil.GetServerLuaFilePathName(s.Sheet.Name)
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(file)

	moduleName := fmt.Sprintf("Schema_%s", strings.ReplaceAll(sheetName, "@", ""))
	//----------------------------字段注释
	write := bufio.NewWriter(file)
	util.WriteStringToFile(write, fmt.Sprintf("---@module %v\n", moduleName))
	for i, cellValue := range s.FieldInfos {
		if i == 0 || i > 1 && !cellValue.SupportLua() {
			continue
		}
		key := strings.TrimSpace(cellValue.Name)
		contentType := strings.TrimSpace(cellValue.FieldType)
		if len(key) > 0 {
			util.WriteStringToFile(write, fmt.Sprintf("---@field %v %v  @%v \n", key, contentType, strings.ReplaceAll(cellValue.Comment, "\n", " ")))
		}
	}

	//----------------------------数据
	util.WriteStringToFile(write, fmt.Sprintf("local %s=\n{\n", moduleName))
	isDirty := false
	for i, cellValue := range s.CellInfos {
		if cellValue == nil || len(cellValue) < 3 {
			continue
		}
		if isDirty {
			util.WriteStringToFile(write, "},\n")
		}
		util.WriteStringToFile(write, fmt.Sprintf("\t[%v] = {", cellValue[1].Content))
		for index, item := range cellValue {
			if index == 0 || index > 1 && !item.FieldInfo.SupportLua() {
				continue
			}
			util.WriteStringToFile(write, fmt.Sprintf("%v = ", item.FieldInfo.Name))
			str, err := exportLua.GetLuaCellValue(item.Content, item.FieldInfo.FieldType)
			if err != nil {
				tip := sheetName + " 行:" + strconv.Itoa(i)
				panic(err.Error() + tip)
			}
			if index != len(cellValue)-1 {
				util.WriteStringToFile(write, fmt.Sprintf("%s, ", str))
			} else {
				util.WriteStringToFile(write, str)
			}
		}
		isDirty = true
	}
	if isDirty {
		util.WriteStringToFile(write, "}\n")
	}

	util.WriteStringToFile(write, fmt.Sprintf("}\n\n---@type table<number,table>\nreturn %v", moduleName))

	err = write.Flush()
	if err != nil {
		return
	}
}

func (s *SheetInfoNormal) ExportCSV() {
	sheetName := s.Sheet.Name
	filePath := fileUtil.GetServerCSVFilePathName(sheetName)
	fileUtil.CreateFile(filePath)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(file)

	write := bufio.NewWriter(file)
	fieldNameStr := ""
	fieldTypeStr := ""
	splitStr := ""
	count := len(s.FieldInfos)
	for j, fieldInfo := range s.FieldInfos {
		if j == 0 || j > 1 && !fieldInfo.SupportCSV() {
			continue
		}
		if j == count-1 {
			splitStr = ""
		} else {
			splitStr = ","
		}
		fieldNameStr = fieldNameStr + fieldInfo.Name + splitStr
		fieldTypeStr = fieldTypeStr + fieldInfo.FieldType + splitStr
	}
	util.WriteStringToFile(write, fieldNameStr+"\n")
	util.WriteStringToFile(write, fieldTypeStr+"\n")

	str := ""
	for _, cellInfo := range s.CellInfos {
		str = ""
		for j, item := range cellInfo {
			if j == 0 || j > 1 && !item.FieldInfo.SupportCSV() {
				continue
			}

			cellValue := strings.TrimSpace(item.Content)
			if len(cellValue) == 0 {
				defaultValue := strings.TrimSpace(item.FieldInfo.DefaultValue)
				cellValue = defaultValue
			}
			if j == count-1 {
				str = str + cellValue
			} else {
				str = str + cellValue + ","
			}
		}
		if str != "" {
			util.WriteStringToFile(write, str+"\n")
		}
	}
	write.Flush()
}
