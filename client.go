package http

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

func (c *Client) Get(url string, data interface{}, opts ...ParamsOption) (*http.Response, error) {
	if c.http == nil {
		c.http = http.DefaultClient
	}
	url = url + "?"
	for i := 0; i < len(opts); i++ {
		url = url + opts[i]()
		if i != len(opts)-1 {
			url = url + "&"
		}
	}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	for _, cookie := range c.Cookie {
		req.AddCookie(cookie)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Post(url string, data interface{}) (*http.Response, error) {
	if c.http == nil {
		c.http = http.DefaultClient
	}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	for _, cookie := range c.Cookie {
		req.AddCookie(cookie)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) PostForm(url string, formData *FormData) (*http.Response, error) {
	if c.http == nil {
		c.http = http.DefaultClient
	}
	req, err := http.NewRequest("POST", url, formData.Buf)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", formData.Write.FormDataContentType())
	for _, cookie := range c.Cookie {
		req.AddCookie(cookie)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) AddCookie(cookie []*http.Cookie) {
	c.Cookie = cookie
}
