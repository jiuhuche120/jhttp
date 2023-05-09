package jhttp

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
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
	client := NewClient(
		AddHeader("Accept", "application/vnd.github.v3+json"),
		AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"),
		AddHeader("Authorization", "Bearer <YOUR-TOKEN>"),
	)
	jsonParams := NewJsonParams(
		AddJsonParam("description", "Example of a gist"),
		AddJsonParam("public", false),
		AddJsonParam("files", "{\"README.md\":{\"content\":\"Hello World\"}}"),
	)
	resp, err := client.Post("https://api.github.com/gists", jsonParams)
	require.Nil(t, err)
	require.Equal(t, true, resp.Contains("Bad credentials"))
}

func TestWebsocket(t *testing.T) {
	client := NewClient()
	ws, resp, err := client.WebSocket("ws://121.40.165.18:8800")
	defer resp.Body.Close()
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			panic(err)
		}
	}(ws)
	require.Nil(t, err)
	_, msg, err := ws.ReadMessage()
	require.Nil(t, err)
	require.Contains(t, string(msg), "服务端主动向你推送")
}
