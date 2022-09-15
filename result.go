package jhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	ReadSize    = 1024 * 1024       //1 MB
	MaxReadSize = 1024 * 1024 * 100 //100 MB
)

// SetReadSize set the slice size to read from the response body
func SetReadSize(size int) {
	ReadSize = size
}

// SetMaxReadSize set the max read size to read from the response body
func SetMaxReadSize(size int) {
	MaxReadSize = size
}

type Result struct {
	resp  *http.Response
	cache []byte
	lock  sync.Mutex
}

// NewResult returns a Result with the given response
func NewResult(resp *http.Response) *Result {
	return &Result{
		resp:  resp,
		cache: nil,
		lock:  sync.Mutex{},
	}
}

// Body read http body from the Result and cache the body
func (result *Result) Body() ([]byte, error) {
	result.lock.Lock()
	defer result.lock.Unlock()
	if result.cache != nil {
		return result.cache, nil
	}
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
	result.cache = data
	return data, nil
}

// JsonUnmarshal json unmarshal the result body into the given type
func (result *Result) JsonUnmarshal(typ interface{}) error {
	body, err := result.Body()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, typ)
	if err != nil {
		return err
	}
	return nil
}

// Header return Result header
func (result *Result) Header() http.Header {
	return result.resp.Header
}

// Cookies return Result cookies
func (result *Result) Cookies() []*http.Cookie {
	return result.resp.Cookies()
}

// StatusCode return Result statusCode
func (result *Result) StatusCode() int {
	return result.resp.StatusCode
}

// Status return Result status
func (result *Result) Status() string {
	return result.resp.Status
}

// IsSuccess return Result success or fail
func (result *Result) IsSuccess() bool {
	return result.StatusCode() == http.StatusOK
}

// Contains return whether Result contains str
func (result *Result) Contains(str string) bool {
	body, err := result.Body()
	if err != nil {
		return false
	}
	if strings.Contains(string(body), str) {
		return true
	}
	return false
}

// Equal return whether Result equal str
func (result *Result) Equal(str string) bool {
	body, err := result.Body()
	if err != nil {
		return false
	}
	if string(body) == str {
		return true
	}
	return false
}
