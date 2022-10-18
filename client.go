package jhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ClientOption = func(*Client)
type ParamsOption = func() string
type Client struct {
	http      *http.Client
	websocket *websocket.Dialer
	header    map[string]string
	cookie    []*http.Cookie
	retry     int
}

// NewClient returns a new Client with ClientOption.
func NewClient(opts ...ClientOption) *Client {
	client := &Client{http: http.DefaultClient, websocket: websocket.DefaultDialer, header: map[string]string{}, retry: 0}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// AddHeader add a header to the client.
func AddHeader(key, value string) ClientOption {
	return func(client *Client) {
		client.header[key] = value
	}
}

// SetTimeout set the timeout for the client.
func SetTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.http.Timeout = timeout
	}
}

// SetRetry set the number of retry for the client. Default is 0.
func SetRetry(retry int) ClientOption {
	return func(client *Client) {
		client.retry = retry
	}
}

// AddParams set the url parameters for the client.
func AddParams(key, value string) ParamsOption {
	return func() string {
		return key + "=" + value
	}
}

// AddCookie set the cookie for the client.
func (c *Client) AddCookie(cookie []*http.Cookie) {
	c.cookie = cookie
}

// Get send a GET to the specified URL.
func (c *Client) Get(url string, data interface{}, opts ...ParamsOption) (*Result, error) {
	url = url + "?"
	for i := 0; i < len(opts); i++ {
		url = url + opts[i]()
		if i != len(opts)-1 {
			url = url + "&"
		}
	}
	return c.doReq(url, "GET", data)
}

// Post send a POST to the specified URL.
func (c *Client) Post(url string, data interface{}) (*Result, error) {
	return c.doReq(url, "POST", data)
}

func (c *Client) WebSocket(url string) (*websocket.Conn, *http.Response, error) {
	header := make(http.Header)
	for k, v := range c.header {
		header.Set(k, v)
	}
	return c.websocket.Dial(url, header)
}

// doReq send appropriate request to the specified URL.
func (c *Client) doReq(url string, reqType string, data interface{}) (*Result, error) {
	switch data.(type) {
	case FormData:
		return c.doForm(url, reqType, data.(FormData))
	case []byte:
		return c.doBytes(url, reqType, data.([]byte))
	case string:
		return c.doString(url, reqType, data.(string))
	default:
		data, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return c.doBytes(url, reqType, data)
	}
}

// doBytes send []byte data to the specified URL.
func (c *Client) doBytes(url string, reqType string, data []byte) (*Result, error) {
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// doString send string data to the specified URL.
func (c *Client) doString(url string, reqType string, data string) (*Result, error) {
	req, err := http.NewRequest(reqType, url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// doForm send FormData to the specified URL.
func (c *Client) doForm(url string, reqType string, formData FormData) (*Result, error) {
	req, err := http.NewRequest(reqType, url, formData.buf)
	if err != nil {
		return nil, err
	}
	// set Form Content-Type
	req.Header.Set("Content-Type", formData.writer.FormDataContentType())
	return c.do(req)
}

// do send an HTTP request and returns an HTTP response. if the request is failed, the client will retry the request until the number of retry.
func (c *Client) do(req *http.Request) (*Result, error) {
	var resp *http.Response
	var err error
	if c.http == nil {
		c.http = http.DefaultClient
	}
	// set http header
	for k, v := range c.header {
		req.Header.Set(k, v)
	}
	// set http cookie
	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}
	for i := 0; i < c.retry+1; i++ {
		resp, err = c.http.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	result, err := NewResult(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}
