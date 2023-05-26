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

func NewClient(opts ...ClientOption) *Client {
	client := &Client{http: http.DefaultClient, websocket: websocket.DefaultDialer, header: map[string]string{}, retry: 0}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func AddHeader(key, value string) ClientOption {
	return func(client *Client) {
		client.header[key] = value
	}
}

func SetTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.http.Timeout = timeout
	}
}

func SetRetry(retry int) ClientOption {
	return func(client *Client) {
		client.retry = retry
	}
}

func AddParams(key, value string) ParamsOption {
	return func() string {
		return key + "=" + value
	}
}

func (c *Client) AddCookie(cookie []*http.Cookie) {
	c.cookie = cookie
}

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

func (c *Client) doReq(url string, reqType string, data interface{}) (*Result, error) {
	switch v := data.(type) {
	case FormData:
		return c.doForm(url, reqType, v)
	case []byte:
		return c.doBytes(url, reqType, v)
	case string:
		return c.doString(url, reqType, v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return c.doBytes(url, reqType, data)
	}
}

func (c *Client) doBytes(url string, reqType string, data []byte) (*Result, error) {
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) doString(url string, reqType string, data string) (*Result, error) {
	req, err := http.NewRequest(reqType, url, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) doForm(url string, reqType string, formData FormData) (*Result, error) {
	req, err := http.NewRequest(reqType, url, formData.buf)
	if err != nil {
		return nil, err
	}
	// set Form Content-Type
	req.Header.Set("Content-Type", formData.writer.FormDataContentType())
	return c.do(req)
}

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
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			} else {
				resp.Body.Close()
			}
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

func (c *Client) GetHeader(key string) string {
	return c.header[key]
}
