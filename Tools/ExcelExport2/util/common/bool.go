package common

import (
	"strings"
)

func ParseBool(fieldType string, value string) (result interface{}) {
	switch fieldType {
	case "bool":
		str := strings.TrimSpace(value)
		if str == "1" {
			result = true
		} else if str == "0" {
			result = false
		} else {
			panic("无法转成bool类型: " + str)
		}
		break
	case "bool[]":
		result = SlipBool(value)
		break
	case "bool[][]":
		result = SlipBool2D(value)
		break
	}
	return
}

func SlipBool(s string) (result []bool) {
	if len(s) == 0 {
		return
	}
	strs := strings.Split(s, "*")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "1" {
			result = append(result, true)
		} else if str == "0" {
			result = append(result, false)
		} else {
			panic("无法转成bool类型: " + str)
		}
	}
	return
}

func SlipBool2D(s string) (result [][]bool) {
	if len(s) == 0 {
		return
	}
	arr := strings.Split(s, "|")
	for _, s := range arr {
		result = append(result, SlipBool(s))
	}
	return
}

func WriteBool(stream *Stream, value interface{}) {
	switch value.(type) {
	case bool:
		stream.TryWriteBool(value.(bool))
		break
	case []bool:
		boolArray := value.([]bool)
		if len(boolArray) > 0 {
			stream.WriteUInt32(uint32(len(boolArray)))
			for i := range boolArray {
				stream.TryWriteBool(boolArray[i])
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	case [][]bool:
		boolArray2D := value.([][]bool)
		if len(boolArray2D) > 0 {
			stream.WriteUInt32(uint32(len(boolArray2D)))
			for _, item := range boolArray2D {
				stream.WriteUInt32(uint32(len(item)))
				for i := range item {
					stream.TryWriteBool(item[i])
				}
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	}
}
