package tool

import "strings"

func CanExportLua(exportType string, isClient bool) bool {
	if len(exportType) == 0 || isClient {
		return false
	}
	return strings.Contains(exportType, "S")
}

func CanExportBin(exportType string, isClient bool) bool {
	if len(exportType) == 0 {
		return false
	}
	return isClient
}

func CanExportCS(exportType string, isClient bool) bool {
	if len(exportType) == 0 {
		return false
	}
	return isClient
}

func CanExportSCSV(exportType string) bool {
	if len(exportType) == 0 {
		return false
	}
	return strings.Contains(exportType, "S")
}
