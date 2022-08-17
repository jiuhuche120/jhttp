package jhttp

import "net/url"

/*
xfrom.go is used to post x-www-form-urlencoded requests
*/

type XFormOption = func(*url.Values)

func AddXFormParams(key, value string) XFormOption {
	return func(values *url.Values) {
		values.Set(key, value)
	}
}

func NewXForm(opts ...XFormOption) string {
	data := url.Values{}
	for _, opt := range opts {
		opt(&data)
	}
	return data.Encode()
}