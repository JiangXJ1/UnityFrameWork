package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Stream struct {
	buf bytes.Buffer
}

func (s *Stream) Len() int {
	return s.buf.Len()
}
func (s *Stream) Buffer() *bytes.Buffer {
	return &s.buf
}
func (s *Stream) WriteRawBytes(b []byte) {
	s.buf.Write(b)
}
func (s *Stream) WriteFile(outfile string) error {
	// 自动创建目录
	os.MkdirAll(filepath.Dir(outfile), 0755)

	err := ioutil.WriteFile(outfile, s.buf.Bytes(), 0666)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	return nil
}
func (s *Stream) WriteInt32(v int32) {
	binary.Write(&s.buf, binary.LittleEndian, v)
}

func (s *Stream) WriteUInt32(v uint32) {
	binary.Write(&s.buf, binary.LittleEndian, v)
}

func (s *Stream) WriteFloat(v float32) {
	binary.Write(&s.buf, binary.LittleEndian, v)
}

func (s *Stream) WriteString(v string) {
	rawStr := []byte(v)

	binary.Write(&s.buf, binary.LittleEndian, int32(len(rawStr)))
	binary.Write(&s.buf, binary.LittleEndian, rawStr)
}

func (s *Stream) WriteBytes(v []byte) {

	binary.Write(&s.buf, binary.LittleEndian, int32(len(v)))

	binary.Write(&s.buf, binary.LittleEndian, v)
}

func (s *Stream) SplitArray(str string) []string {
	strArray := strings.Split(str, "*")
	return strArray
}

func (s *Stream) SplitArray2D(str string) [][]string {
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

func (s *Stream) GetArray(value string) []string {
	array := s.SplitArray(value)
	len := len(array)
	s.WriteUInt32(uint32(len))
	return array
}

func (s *Stream) GetArray2D(value string) []string {
	strArray := strings.Split(value, "|")
	len := len(strArray)
	s.WriteUInt32(uint32(len))
	return strArray
}

func (s *Stream) TryWriteInt32(value int32) {
	s.WriteInt32(value)
}
func (s *Stream) TryWriteUInt32(value uint32) {
	s.WriteUInt32(value)
}
func (s *Stream) TryWriteFloat(value float32) {
	s.WriteFloat(value)
}

func (s *Stream) TryWriteBool(value bool) {
	if value {
		s.WriteUInt32(1)
	} else {
		s.WriteUInt32(0)
	}
}

func (s *Stream) WriteNodeValue(strType string, value interface{}) {
	switch value.(type) {
	case bool, []bool, [][]bool:
		WriteBool(s, value)
		break
	case int32, []int32, [][]int32:
		WriteInt32(s, value)
		break
	case uint32, []uint32, [][]uint32:
		WriteUint32(s, value)
		break
	case float32, []float32, [][]float32:
		WriteFloat(s, value)
		break
	case string, []string, [][]string:
		WriteString(s, value)
		break
	default:
		fmt.Printf("type:%v\n", reflect.TypeOf(value))
		panic("unsupport type:" + strType)
	}
}

func NewStream() *Stream {
	return &Stream{}
}
