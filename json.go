package jhttp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type JsonOption func(jsonStruct *JsonStruct)
type JsonStruct struct {
	jsonMp map[string]interface{}
}

func NewJsonParams(opts ...JsonOption) []byte {
	var jsonStruct JsonStruct
	jsonStruct.jsonMp = make(map[string]interface{})
	for _, opt := range opts {
		opt(&jsonStruct)
	}

	str := strings.Builder{}
	str.WriteString("{")
	length := len(jsonStruct.jsonMp)
	index := 1
	for param, value := range jsonStruct.jsonMp {
		switch v := value.(type) {
		case string:
			// escape `"`
			v = strconv.Quote(v)
			v = v[1 : len(v)-1]
			str.WriteString(fmt.Sprintf("\"%v\":\"%v\"", param, v))
		case bool, int, int8, int16, int32, int64, float32, float64:
			str.WriteString(fmt.Sprintf("\"%v\":%v", param, v))
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return nil
			}
			str.WriteString(fmt.Sprintf("\"%v\":%v", param, string(data)))
		}
		if index < length {
			index++
			str.WriteString(",")
		}
	}
	str.WriteString("}")
	return []byte(str.String())
}

func AddJsonParam(key string, value interface{}) JsonOption {
	return func(jsonStruct *JsonStruct) {
		jsonStruct.jsonMp[key] = value
	}
}
