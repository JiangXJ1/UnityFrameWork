package main

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil"
	"OneExcel/src/fileUtil"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	//if excelUtil.IsExcelClosed() != true {
	//	return
	//}
	ExportAllExcel()
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}

func ExportAllExcel() {
	timeNow := time.Now().UnixNano()
	fmt.Println("Exporting...")

	defer func() {
		if e := recover(); e != nil {
			fmt.Println("错误详细信息:\n\n")
			PrintStack()
		}
		fmt.Println(src.ErrorStr)
		src.ErrorStr = ""
		fmt.Printf("Press any key to exit...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
	}()

	src.ClearPath()
	ExportAllExcelSheetData(excelUtil.GetAllData())

	timeNow = time.Now().UnixNano() - timeNow
	timeNow = timeNow / 1e6
	fmt.Printf("\nTotal time : %d ms. \n", timeNow)
	fmt.Println("\nExport Success...")

	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

func ExportAllExcelSheetData(sheetDatas []src.SheetData) {
	filePath := src.OUTPUT_DIR_CS + "SchemaCreator.cs"
	if fileUtil.FileExist(filePath) {
		fileUtil.DeleteFile(filePath)
	}

	for i := range sheetDatas {
		sheetData := sheetDatas[i]
		ExportExcelSheetData(sheetData)
	}
}

func ExportExcelSheetData(sheetData src.SheetData) {
	src.CurrentSheetData = sheetData
	excelUtil.WriteData(sheetData)
}
