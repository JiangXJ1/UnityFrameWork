package export_csharp

import (
	"OneExcel/src"
	"OneExcel/src/excelUtil/common"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Stream struct {
	buf bytes.Buffer
}

func (self *Stream) Len() int {
	return self.buf.Len()
}
func (self *Stream) Buffer() *bytes.Buffer {
	return &self.buf
}
func (self *Stream) WriteRawBytes(b []byte) {
	self.buf.Write(b)
}
func (self *Stream) WriteFile(outfile string) error {
	// 自动创建目录
	os.MkdirAll(filepath.Dir(outfile), 0755)

	err := ioutil.WriteFile(outfile, self.buf.Bytes(), 0666)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	return nil
}
func (self *Stream) WriteInt32(v int32) {

	binary.Write(&self.buf, binary.LittleEndian, v)
}

func (self *Stream) WriteUInt32(v uint32) {

	binary.Write(&self.buf, binary.LittleEndian, v)
}

func (self *Stream) WriteFloat(v float32) {
	binary.Write(&self.buf, binary.LittleEndian, v)
}

func (self *Stream) WriteString(v string) {
	rawStr := []byte(v)

	binary.Write(&self.buf, binary.LittleEndian, int32(len(rawStr)))

	binary.Write(&self.buf, binary.LittleEndian, rawStr)
}

func (self *Stream) WriteBytes(v []byte) {

	binary.Write(&self.buf, binary.LittleEndian, int32(len(v)))

	binary.Write(&self.buf, binary.LittleEndian, v)
}

func (self *Stream) SplitArray(str string) []string {
	strArray := strings.Split(str, "*")
	return strArray
}

func (self *Stream) SplitArray2D(str string) [][]string {
	var contentList [][]string
	strArray := strings.Split(str, "|")
	for i := range strArray {
		var cell []string
		cellArray := strings.Split(strArray[i], "*")
		for j := range cellArray {
			cell = append(cell, cellArray[j])
		}
		contentList = append(contentList, cell)
	}
	return contentList
}

func (self *Stream) GetArray(value string) []string {
	array := self.SplitArray(value)
	len := len(array)
	self.WriteUInt32(uint32(len))
	return array
}

func (self *Stream) GetArray2D(value string) []string {
	strArray := strings.Split(value, "|")
	len := len(strArray)
	self.WriteUInt32(uint32(len))
	return strArray
}

func (self *Stream) TryWriteInt32(value string) {
	v, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		src.IPanicEx(" ParseInt error! " + value)
	}
	self.WriteInt32(int32(v))
}
func (self *Stream) TryWriteUInt32(value string) {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		src.IPanicEx("ParseUint error! " + value)
	}
	self.WriteUInt32(uint32(v))
}
func (self *Stream) TryWriteFloat(value string) {
	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		src.IPanicEx("ParseFloat error! " + value)
	}
	self.WriteFloat(float32(v))
}

func (self *Stream) TryWriteBool(value string) {
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil || (v != 0 && v != 1) {
		src.IPanicEx("ParseFloat error!bool 必须填1或者0 " + value)
	}
	self.WriteUInt32(uint32(v))
}

func (self *Stream) WriteNodeValue(strType string, value string) {
	ft := common.GetFieldType(strType)
	switch ft {
	case common.FieldType_Bool:
		self.TryWriteBool(value)
		break
	case common.FieldType_BoolArray:
		if len(value) > 0 {
			array := self.GetArray(value)
			for i := range array {
				self.TryWriteBool(array[i])
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_BoolArray2D:
		if len(value) > 0 {
			arrayRoot := self.GetArray2D(value)
			for i := range arrayRoot {
				array := self.GetArray(arrayRoot[i])
				for i := range array {
					self.TryWriteBool(array[i])
				}
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_Int32:
		self.TryWriteInt32(value)
		break
	case common.FieldType_Int32Array:
		if len(value) > 0 {
			array := self.GetArray(value)
			for i := range array {
				self.TryWriteInt32(array[i])
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_Int32Array2D:
		if len(value) > 0 {
			arrayRoot := self.GetArray2D(value)
			for i := range arrayRoot {
				array := self.GetArray(arrayRoot[i])
				for i := range array {
					self.TryWriteInt32(array[i])
				}
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_UInt32:
		self.TryWriteUInt32(value)
		break
	case common.FieldType_UInt32Array:
		if len(value) > 0 {
			array := self.GetArray(value)
			for i := range array {
				self.TryWriteUInt32(array[i])
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_UInt32Array2D:
		if len(value) > 0 {
			arrayRoot := self.GetArray2D(value)
			for i := range arrayRoot {
				array := self.GetArray(arrayRoot[i])
				for i := range array {
					self.TryWriteUInt32(array[i])
				}
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_Float:
		self.TryWriteFloat(value)
		break
	case common.FieldType_FloatArray:
		if len(value) > 0 {
			array := self.GetArray(value)
			for i := range array {
				self.TryWriteFloat(array[i])
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_FloatArray2D:
		if len(value) > 0 {
			arrayRoot := self.GetArray2D(value)
			for i := range arrayRoot {
				array := self.GetArray(arrayRoot[i])
				for i := range array {
					self.TryWriteFloat(array[i])
				}
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_String:
		self.WriteString(value)
		break
	case common.FieldType_StringArray:
		if len(value) > 0 {
			array := self.GetArray(value)
			for i := range array {
				self.WriteString(array[i])
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	case common.FieldType_StringArray2D:
		if len(value) > 0 {
			arrayRoot := self.GetArray2D(value)
			for i := range arrayRoot {
				array := self.GetArray(arrayRoot[i])
				for i := range array {
					self.WriteString(array[i])
				}
			}
		} else {
			self.WriteUInt32(uint32(0))
		}
		break
	default:
		src.IPanicEx("unsupport type:" + strType)
		panic("unsupport type:" + strType)
	}
}

func NewStream() *Stream {
	return &Stream{}
}
