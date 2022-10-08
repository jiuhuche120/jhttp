package jhttp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewXFormParams(t *testing.T) {
	xFormParams := NewXFormParams(
		AddXFormParams("k1", "v1"),
		AddXFormParams("k2", "v2"),
	)
	require.Equal(t, "k1=v1&k2=v2", xFormParams)
}
