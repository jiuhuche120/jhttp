package jhttp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		AddHeader("Accept", "application/vnd.github.v3+json"),
		SetTimeout(time.Second*10),
		SetRetry(3),
	)
	require.NotNil(t, client)
	require.Equal(t, "application/vnd.github.v3+json", client.header["Accept"])
}

func TestGet(t *testing.T) {
	client := NewClient(
		AddHeader("Accept", "application/vnd.github.v3+json"),
		AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"),
	)
	opts := []ParamsOption{
		AddParams("wd", "github"),
	}
	resp, err := client.Get("https://www.baidu.com/s", nil, opts...)
	require.Nil(t, err)
	require.Equal(t, true, resp.Contains("GitHub · Build software better, together."))
}

func TestPost(t *testing.T) {

}
