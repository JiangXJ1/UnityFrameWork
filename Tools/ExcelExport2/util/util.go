package util

import (
	"bufio"
	"strings"
)

func IsExcelFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".xlsx") || strings.HasSuffix(fileName, ".xls") || strings.HasSuffix(fileName, ".xlsm")
}

func WriteStringToFile(writer *bufio.Writer, content string) {
	_, err := writer.WriteString(content)
	if err != nil {
		panic(err)
	}
}
