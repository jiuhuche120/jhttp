package jhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ClientOption = func(*Client)
type ParamsOption = func() string
type Client struct {
	http   *http.Client
	Header map[string]string
	Cookie []*http.Cookie
}

func NewClient(opts ...ClientOption) *Client {
	client := &Client{http: http.DefaultClient, Header: map[string]string{}}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func AddHeader(key, value string) ClientOption {
	return func(client *Client) {
		client.Header[key] = value
	}
}

func AddParams(key, value string) ParamsOption {
	return func() string {
		return key + "=" + value
	}
}

func (c *Client) AddCookie(cookie []*http.Cookie) {
	c.Cookie = cookie
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

func (c *Client) doReq(url string, reqType string, data interface{}) (*Result, error) {
	switch data.(type) {
	case FormData:
		return c.doForm(url, reqType, data.(*FormData))
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

func (c *Client) doForm(url string, reqType string, formData *FormData) (*Result, error) {
	req, err := http.NewRequest(reqType, url, formData.Buf)
	if err != nil {
		return nil, err
	}
	// set Form Content-Type
	req.Header.Set("Content-Type", formData.Write.FormDataContentType())
	return c.do(req)
}

func (c *Client) do(req *http.Request) (*Result, error) {
	if c.http == nil {
		c.http = http.DefaultClient
	}
	// set http header
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	// set http cookie
	for _, cookie := range c.Cookie {
		req.AddCookie(cookie)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return &Result{*resp, nil}, nil
}
