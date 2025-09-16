package common

import (
	"strconv"
	"strings"
)

func ParseUInt32(fieldType string, value string) (result interface{}) {
	switch fieldType {
	case "uint":
		if len(value) == 0 {
			panic("parse uint32 error")
		} else {
			target, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				panic("parse uint32 error" + value)
			}
			result = uint32(target)
		}
		break
	case "uint[]":
		result = SlipUInt(value)
		break
	case "uint[][]":
		result = SlipUInt2D(value)
		break
	}
	return
}

func SlipUInt(s string) (result []uint32) {
	if len(s) == 0 {
		return
	}
	strs := strings.Split(s, "*")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			result = append(result, uint32(0))
			continue
		}
		//尝试将字符串转换为uint32
		f, err := strconv.ParseUint(str, 10, 32)
		if err == nil {
			result = append(result, uint32(f))
		} else {
			panic("parse uint32 error")
		}
	}
	return
}

func SlipUInt2D(s string) (result [][]uint32) {
	if len(s) == 0 {
		return
	}

	arr := strings.Split(s, "|")
	for _, s := range arr {
		result = append(result, SlipUInt(s))
	}
	return
}

func WriteUint32(stream *Stream, value interface{}) {
	switch value.(type) {
	case uint32:
		stream.TryWriteUInt32(value.(uint32))
		break
	case []uint32:
		uint32Array := value.([]uint32)
		if len(uint32Array) > 0 {
			stream.WriteUInt32(uint32(len(uint32Array)))
			for i := range uint32Array {
				stream.TryWriteUInt32(uint32Array[i])
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	case [][]uint32:
		uint32Array2D := value.([][]uint32)
		if len(uint32Array2D) > 0 {
			stream.WriteUInt32(uint32(len(uint32Array2D)))
			for _, item := range uint32Array2D {
				stream.WriteUInt32(uint32(len(item)))
				for i := range item {
					stream.TryWriteUInt32(item[i])
				}
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	}
}
