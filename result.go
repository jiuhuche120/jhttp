package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ReadSize    = 1024 * 1024       //1 MB
	MaxReadSize = 1024 * 1024 * 100 //100 MB
)

func SetReadSize(size int) {
	ReadSize = size
}

func SetMaxReadSize(size int) {
	MaxReadSize = size
}

type Result struct {
	resp http.Response
}

func (result *Result) Body() ([]byte, error) {
	readSlice := make([]byte, ReadSize)
	var data []byte
	for size, err := result.resp.Body.Read(readSlice); size != 0; size, err = result.resp.Body.Read(readSlice) {
		data = append(data, readSlice[:size]...)
		if len(data) > MaxReadSize {
			return nil, fmt.Errorf("too many bytes to read")
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
	}
	return data, nil
}

func (result *Result) JsonUnmarshal(data interface{}) error {
	body, err := result.Body()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

func (result *Result) Header() http.Header {
	return result.resp.Header
}

func (result *Result) Cookies() []*http.Cookie {
	return result.resp.Cookies()
}

func (result *Result) StatusCode() int {
	return result.resp.StatusCode
}

func (result *Result) Status() string {
	return result.resp.Status
}

func (result *Result) IsSuccess() bool {
	if result.StatusCode() == http.StatusOK {
		return true
	}
	return false
}

func (result *Result) Contains(key string) bool {
	body, err := result.Body()
	if err != nil {
		return false
	}
	if strings.Contains(string(body), key) {
		return true
	}
	return false
}

func (result *Result) Equal(key string) bool {
	body, err := result.Body()
	if err != nil {
		return false
	}
	if string(body) == key {
		return true
	}
	return false
}
