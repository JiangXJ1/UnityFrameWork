package provider

import "fmt"

// ExportError 自定义导出错误
type ExportError struct {
	FileName  string // Excel文件名
	SheetName string // 表名
	FieldName string // 字段名
	RowNum    int    // 行号
	Message   string // 错误信息
}

func (e *ExportError) Error() string {
	return fmt.Sprintf("\033[31m导出错误 - 【文件:%s 页签:%s 字段:%s 行:%d 错误详情:%s】\033[0m",
		e.FileName, e.SheetName, e.FieldName, e.RowNum, e.Message)
}

func NewExportError(fileName, sheetName, fieldName string, rowNum int, message string) *ExportError {
	return &ExportError{
		FileName:  fileName,
		SheetName: sheetName,
		FieldName: fieldName,
		RowNum:    rowNum,
		Message:   message,
	}
}
