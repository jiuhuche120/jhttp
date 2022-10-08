package jhttp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFormParams(t *testing.T) {
	formParams, err := NewFormParams(
		AddFormParams("username", "username", Text),
		AddFormParams("password", "password", Text),
	)
	require.Nil(t, err)
	fmt.Println(formParams.buf.String())
}
