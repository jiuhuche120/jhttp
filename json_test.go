package jhttp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MyStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func TestNewJsonParams(t *testing.T) {
	var arr []MyStruct
	arr = append(arr, MyStruct{Key: "k1", Value: "v1"})
	arr = append(arr, MyStruct{Key: "k2", Value: "v2"})
	arr = append(arr, MyStruct{Key: "k3", Value: "v3"})
	data := NewJsonParams(
		AddJsonParam("k1", "v1"),
		AddJsonParam("k2", arr),
	)
	require.Equal(t, "{\"k1\":\"v1\",\"k2\":[{\"key\":\"k1\",\"value\":\"v1\"},{\"key\":\"k2\",\"value\":\"v2\"},{\"key\":\"k3\",\"value\":\"v3\"}]}", string(data))
}
