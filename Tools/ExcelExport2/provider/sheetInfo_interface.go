package provider

import "github.com/tealeg/xlsx"

type SheetInfoBase struct {
	FileName   string
	Sheet      *xlsx.Sheet
	FieldInfos []FieldInfo
	CellInfos  [][]CellInfo
	CanExportC bool
	CanExportS bool
}

type SheetInfoInterface interface {
	Read(task *Task)
	Export()
}
