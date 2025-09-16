package common

import (
	"strings"
)

func ParseString(fieldType string, value string) (result interface{}) {
	switch fieldType {
	case "string":
		result = value
		break
	case "string[]":
		result = SlipString(value)
		break
	case "string[][]":
		result = SlipString2D(value)
		break
	}
	return
}

func SlipString(s string) (result []string) {
	if len(s) == 0 {
		return
	}
	result = strings.Split(s, "*")
	return
}

func SlipString2D(s string) (result [][]string) {
	if len(s) == 0 {
		return
	}
	arr := strings.Split(s, "|")
	for _, s := range arr {
		result = append(result, SlipString(s))
	}
	return
}

func WriteString(stream *Stream, value interface{}) {
	switch value.(type) {
	case string:
		stream.WriteString(value.(string))
		break
	case []string:
		stringArray := value.([]string)
		if len(stringArray) > 0 {
			stream.WriteUInt32(uint32(len(stringArray)))
			for i := range stringArray {
				stream.WriteString(stringArray[i])
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	case [][]string:
		stringArray2D := value.([][]string)
		if len(stringArray2D) > 0 {
			stream.WriteUInt32(uint32(len(stringArray2D)))
			for _, item := range stringArray2D {
				stream.WriteUInt32(uint32(len(item)))
				for i := range item {
					stream.WriteString(item[i])
				}
			}
		} else {
			stream.WriteUInt32(uint32(0))
		}
		break
	}
}
