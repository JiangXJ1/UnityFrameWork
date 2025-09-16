package src

import "fmt"

var CurrentRow int = 0
var CurrentColum int = 0
var CurrentSheetName string = ""
var CurrentFiledName string = ""
var CurrentFileName string = ""
var CurrentSheetData SheetData

var ErrorStr string = ""

func IPanic(info string) {
	PrintTableInfo(info)
	panic("")
}

func PrintTableInfo(info string) {
	pointStr := "******************************************************************\n"
	ErrorStr = pointStr
	ErrorStr = ErrorStr + info + "\n\n"
	ErrorStr = ErrorStr + fmt.Sprintf("[%v]--->[%v]--->[第%d行]--->[第%d列] --->[字段名:%s]", CurrentFileName, CurrentSheetName, CurrentRow+1, CurrentColum+1, CurrentFiledName) + "\n"
	ErrorStr = ErrorStr + pointStr
	//fmt.Println("******************************************************************")
	//fmt.Println(info + "\n")
	//_, err := fmt.Printf("[%v]--->[%v]--->[第%d行]--->[第%d列]", CurrentFileName, CurrentSheetName, CurrentRow+1, CurrentColum+1)
	//if err == nil {
	//	fmt.Println()
	//} else {
	//	fmt.Println(err)
	//}
	//fmt.Println("******************************************************************")
}

func IPanicEx(info string) {
	ErrorStr = ""
	ErrorStr = ErrorStr + info + "\n"
	ErrorStr = fmt.Sprintf("[%v]--->[%v]", CurrentSheetData.FileName, CurrentSheetData.SheetName)
	//fmt.Println("******************************************************************")
	//fmt.Println(info + "\n")
	//_, err := fmt.Printf("[%v]--->[%v]", CurrentSheetData.FileName, CurrentSheetData.SheetName)
	//if err == nil {
	//	fmt.Println()
	//} else {
	//	fmt.Println(err)
	//}
	//fmt.Println("******************************************************************")
	panic("")
}
