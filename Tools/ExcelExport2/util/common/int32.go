package common

import (
	"strconv"
	"strings"
)

func ParseInt32(fieldType string, value string) (result interface{}) {
	switch fieldType {
	case "int":
		if len(value) == 0 {
			panic("parse int32 error")
		}
		f, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			panic("parse int32 error")
		}
		result = int32(f)
		break
	case "int[]":
		result = SlipInt(value)
		break
	case "int[][]":
		result = SlipInt2D(value)
		break
	}
	return
}

func SlipInt(s string) (result []int32) {
	if len(s) == 0 {
		return
	}
	strs := strings.Split(s, "*")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			result = append(result, int32(0))
			continue
		}
		//尝试将字符串转换为int32
		f, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			result = append(result, int32(f))
		} else {
			panic("parse int32 error")
		}
	}
	return
}

func SlipInt2D(s string) (result [][]int32) {
	if len(s) == 0 {
		return
	}
	arr := strings.Split(s, "|")
	for _, s := range arr {
		result = append(result, SlipInt(s))
	}
	return
}

func WriteInt32(stream *Stream, value interface{}) {
	switch value.(type) {
	case int32:
		stream.TryWriteInt32(value.(int32))
		break
	case []int32:
		int32Array := value.([]int32)
		if len(int32Array) > 0 {
			stream.WriteUInt32(uint32(len(int32Array)))
			for i := range int32Array {
				stream.TryWriteInt32(int32Array[i])
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	case [][]int32:
		int32Array2D := value.([][]int32)
		if len(int32Array2D) > 0 {
			stream.WriteUInt32(uint32(len(int32Array2D)))
			for _, item := range int32Array2D {
				stream.WriteUInt32(uint32(len(item)))
				for i := range item {
					stream.TryWriteInt32(item[i])
				}
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	default:
		panic("invalid type")
	}
}
