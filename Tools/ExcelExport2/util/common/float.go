package common

import (
	"strconv"
	"strings"
)

func ParseFloat(fieldType string, value string) (result interface{}) {
	switch fieldType {
	case "float":
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			panic("parse float32 error")
		}
		result = float32(f)
		break
	case "float[]":
		result = SlipFloat(value)
		break
	case "float[][]":
		result = SlipFloat2D(value)
		break
	}
	return
}

func SlipFloat(s string) (result []float32) {
	if len(s) == 0 {
		return
	}
	strs := strings.Split(s, "*")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			result = append(result, float32(0))
			continue
		}
		//尝试将字符串转换为float32
		f, err := strconv.ParseFloat(str, 32)
		if err == nil {
			result = append(result, float32(f))
		} else {
			panic("parse float error str:" + str)
		}

	}
	return
}

func SlipFloat2D(s string) (result [][]float32) {
	if len(s) == 0 {
		return
	}
	arr := strings.Split(s, "|")
	for _, s := range arr {
		result = append(result, SlipFloat(s))
	}
	return
}

func WriteFloat(stream *Stream, value interface{}) {
	switch value.(type) {
	case float32:
		stream.TryWriteFloat(value.(float32))
		break
	case []float32:
		float32Array := value.([]float32)
		if len(float32Array) > 0 {
			stream.WriteUInt32(uint32(len(float32Array)))
			for i := range float32Array {
				stream.TryWriteFloat(float32Array[i])
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	case [][]float32:
		float32Array2D := value.([][]float32)
		if len(float32Array2D) > 0 {
			stream.WriteUInt32(uint32(len(float32Array2D)))
			for _, item := range float32Array2D {
				stream.WriteUInt32(uint32(len(item)))
				for i := range item {
					stream.TryWriteFloat(item[i])
				}
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	}
}
