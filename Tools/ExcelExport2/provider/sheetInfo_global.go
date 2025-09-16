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

type SheetInfoGlobal struct {
	SheetInfoBase
	LocalCellInfos []CellInfo
}

func (s *SheetInfoGlobal) Read(task *Task) {
	s.FileName = task.FileName
	s.Sheet = task.Sheet
	//解析字段信息和内容
	s.GetFieldInfo()
}

// 解析字段信息
func (s *SheetInfoGlobal) GetFieldInfo() {
	rowLen := len(s.Sheet.Rows)
	columnCount := rowLen - 5
	s.FieldInfos = make([]FieldInfo, 0, columnCount)
	s.LocalCellInfos = make([]CellInfo, 0, columnCount)
	names := make(map[string]bool)

	for i := 5; i < rowLen; i++ {
		cells := s.Sheet.Rows[i].Cells
		if len(cells) < 5 || len(strings.TrimSpace(cells[3].Value)) == 0 || strings.HasPrefix(strings.TrimSpace(cells[0].Value), "#") {
			continue
		}
		//字段信息
		comment := ""
		fieldName := cells[1].Value
		fieldType := cells[2].Value
		exportType := cells[3].Value
		content := cells[4].Value
		defaultValue := ""
		if len(cells) >= 6 {
			comment = cells[5].Value
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
		s.FieldInfos = append(s.FieldInfos, fieldInfo)
		//内容
		if fieldInfo.Valid {
			cellInfo := CellInfo{
				FieldInfo: fieldInfo,
				Content:   content,
			}
			// 添加错误处理
			func() {
				defer func() {
					if r := recover(); r != nil {
						panic(NewExportError(
							s.FileName,
							s.Sheet.Name,
							fieldName,
							i+1,
							fmt.Sprint(r),
						))
					}
				}()
				cellInfo.TryParse()
			}()
			s.LocalCellInfos = append(s.LocalCellInfos, cellInfo)
		}
	}
	for _, item := range s.FieldInfos {
		if !item.Valid {
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

// Export 导出
func (s *SheetInfoGlobal) Export() {
	if len(s.FieldInfos) == 0 {
		return
	}
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
func (s *SheetInfoGlobal) ExportCsharp() {
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
		if !fieldInfo.Valid || !fieldInfo.SupportClient() {
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
	CsharpGlobalTemplateParam := export.CsharpGlobalTemplateParam{
		ClassName:   fmt.Sprintf("Schema_%s", strings.ReplaceAll(s.Sheet.Name, "@", "")),
		KeyName:     s.FieldInfos[0].Name,
		KeyType:     s.FieldInfos[0].FieldType,
		DataDefine:  strDefine,
		Deserialize: strDeserialize,
	}
	//生成代码
	content := CsharpGlobalTemplateParam.GenerateCsharpTemplate()
	//写入文件
	write := bufio.NewWriter(file)
	util.WriteStringToFile(write, content)
	err = write.Flush()
	if err != nil {
		return
	}
}

// ExportBin 导出二进制文件
func (s *SheetInfoGlobal) ExportBin() {
	fileStream := common.NewStream()
	rowStream := common.NewStream()
	for _, cellValue := range s.LocalCellInfos {
		if !cellValue.FieldInfo.Valid || !cellValue.FieldInfo.SupportClient() {
			continue
		}
		//行
		fieldType := cellValue.FieldInfo.FieldType
		rowStream.WriteNodeValue(fieldType, cellValue.Value)
	}
	fileStream.WriteInt32(int32(0))
	fileStream.WriteInt32(int32(rowStream.Len()))
	fileStream.WriteRawBytes(rowStream.Buffer().Bytes())

	err := fileStream.WriteFile(fileUtil.GetBinFilePathName(s.Sheet.Name))
	if err != nil {
		fmt.Printf("导出二进制文件失败：%s\n", err.Error())
	}
}

func (s *SheetInfoGlobal) ExportLua() {
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

	sheetName := strings.ReplaceAll(s.Sheet.Name, "@", "")
	moduleName := fmt.Sprintf("Schema_%s", sheetName)
	//写入文件
	write := bufio.NewWriter(file)
	util.WriteStringToFile(write, fmt.Sprintf("---@class %v\n", moduleName))

	for i := range s.LocalCellInfos {
		cellValue := s.LocalCellInfos[i]
		if !cellValue.FieldInfo.SupportLua() {
			continue
		}
		key := strings.TrimSpace(cellValue.FieldInfo.Name)
		contentType := strings.TrimSpace(cellValue.FieldInfo.FieldType)
		if len(key) > 0 {
			util.WriteStringToFile(write, fmt.Sprintf("---@field %v %v  @%v \n", key, contentType, strings.ReplaceAll(cellValue.FieldInfo.Comment, "\n", " ")))
		}
	}

	util.WriteStringToFile(write, fmt.Sprintf("local %s = \n{\n", moduleName))

	for i := range s.LocalCellInfos {
		cellValue := s.LocalCellInfos[i]
		if !cellValue.FieldInfo.Valid || !cellValue.FieldInfo.SupportServer() {
			continue
		}
		util.WriteStringToFile(write, fmt.Sprintf("\t%s = ", cellValue.FieldInfo.Name))
		fixCellValue, err := exportLua.GetLuaCellValue(cellValue.Content, cellValue.FieldInfo.FieldType)
		if err != nil {
			tip := sheetName + " 行:" + strconv.Itoa(i)
			panic(err.Error() + tip)
		}
		util.WriteStringToFile(write, fmt.Sprintf("%s,\n", fixCellValue))
	}
	util.WriteStringToFile(write, fmt.Sprintf("}\n\nreturn %v", moduleName))

	err = write.Flush()
	if err != nil {
		return
	}
}

func (s *SheetInfoGlobal) ExportCSV() {
	sheetName := strings.ReplaceAll(s.Sheet.Name, "@", "")
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

	util.WriteStringToFile(write, "KEY,Type,Content\n")
	str := ""

	rowLength := len(s.LocalCellInfos)
	for i, cellValues := range s.LocalCellInfos {
		if !cellValues.FieldInfo.SupportCSV() {
			continue
		}
		key := strings.TrimSpace(cellValues.FieldInfo.Name)
		contentType := strings.TrimSpace(cellValues.FieldInfo.FieldType)
		if i == rowLength-1 {
			str = key + "," + contentType + "," + cellValues.Content
		} else {
			str = key + "," + contentType + "," + cellValues.Content + "\n"
		}
		util.WriteStringToFile(write, str)
	}
	write.Flush()
}
