package cs

import (
	"OneExcel/src"
	"OneExcel/src/fileUtil"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func WriteEnum(sheetData src.SheetData) {
	content := ""
	fileUtil.CreateDir(src.OUTPUT_DIR_CS)
	filePath := src.OUTPUT_DIR_CS + "SchemaCreator.cs"
	if !fileUtil.FileExist(filePath) {
		fileUtil.CreateFile(filePath)
		content = ConfigEnum
	} else {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("read file err:", err.Error())
			return
		}
		content = string(data)

		if strings.Contains(sheetData.SheetName, "Language_") && strings.Contains(content, "Schema_Language") {
			return
		}
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	if !strings.Contains(content, "//SCHEMA_BEGIN") {
		content = ConfigEnum
	}

	name := strings.ReplaceAll(sheetData.SheetName, "@", "")
	if strings.Contains(name, "Language_") {
		name = "Schema_Language"
	}

	strID := name + ","
	content = strings.ReplaceAll(content, "//SCHEMA_ID_BEGIN", "//SCHEMA_ID_BEGIN\n\t\t"+strID)

	str := "case SchemaID." + name + ": return new Schema." + name + "();"
	content = strings.ReplaceAll(content, "//SCHEMA_BEGIN", "//SCHEMA_BEGIN\n\t\t\t\t"+str)

	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
}
