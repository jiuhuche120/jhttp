package jhttp

import "net/url"

type XFormOption = func(*url.Values)

func AddXFormParams(key, value string) XFormOption {
	return func(values *url.Values) {
		values.Set(key, value)
	}
}

func NewXFormParams(opts ...XFormOption) string {
	data := url.Values{}
	for _, opt := range opts {
		opt(&data)
	}
	return data.Encode()
}
