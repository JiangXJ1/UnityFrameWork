package provider

import (
	"demo-go-excel/global"
	"fmt"
	"github.com/tealeg/xlsx"
	"sync/atomic"
)

// 定义任务类型
type Task struct {
	RunInfo  *RunInfo
	Sheet    *xlsx.Sheet
	IsDirty  bool
	FileName string
}

// 执行任务
func (t *Task) Execute() {
	defer func() {
		atomic.AddInt32(&global.CurTaskCount, 1)
		if r := recover(); r != nil {
			if err, ok := r.(*ExportError); ok {
				t.RunInfo.errChan <- err
			} else {
				t.RunInfo.errChan <- NewExportError(t.FileName, t.Sheet.Name, "", 0, fmt.Sprint(r))
			}
		}
	}()

	// 检查是否已有错误发生
	if t.RunInfo.HasError() {
		return
	}

	if len(t.Sheet.Rows) <= global.MIN_ROW || len(t.Sheet.Cols) < global.MIN_COL {
		return
	}

	var sheetInfo SheetInfoInterface
	switch {
	case global.IsGlobal(t.Sheet.Name):
		sheetInfo = new(SheetInfoGlobal)
	case global.IsLanguage(t.Sheet.Name):
		sheetInfo = new(SheetInfoLanguage)
	default:
		sheetInfo = new(SheetInfoNormal)
	}
	sheetInfo.Read(t)
	sheetInfo.Export()
}
