package jhttp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
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
	l := len(jsonStruct.jsonMp)
	index := 1
	for key, value := range jsonStruct.jsonMp {
		switch value.(type) {
		case string:
			str.WriteString(fmt.Sprintf("\"%v\":\"%v\"", key, value))
		default:
			str.WriteString(fmt.Sprintf("\"%v\":%v", key, value))
		}
		if index < l {
			index++
			str.WriteString(",")
		}
	}
	str.WriteString("}")
	value := gjson.Parse(str.String()).Value()
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}

func AddJsonParam[T string | bool | int | int8 | int16 | int32 | int64 | float32 | float64](key string, value T) JsonOption {
	return func(jsonStruct *JsonStruct) {
		jsonStruct.jsonMp[key] = value
	}
}
