package provider

import "strings"

type FieldInfo struct {
	Index        int
	Valid        bool
	ExportType   string
	Name         string
	FieldType    string
	DefaultValue string
	Comment      string
}

func (f *FieldInfo) SupportClient() bool {
	if !f.Valid {
		return false
	}
	return strings.Contains(f.ExportType, "C")
}

func (f *FieldInfo) SupportServer() bool {
	if !f.Valid {
		return false
	}
	return strings.Contains(f.ExportType, "S")
}

func (f *FieldInfo) SupportLua() bool {
	return f.SupportServer()
}

func (f *FieldInfo) SupportCSV() bool {
	return f.SupportServer()
}

func (f *FieldInfo) SupportCSharp() bool {
	return f.SupportClient()
}
