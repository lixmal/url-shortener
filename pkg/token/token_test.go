package token_test

import (
	"github.com/lixmal/url-shortener/pkg/token"
	"github.com/stretchr/testify/assert"
)
import "testing"

func TestNew(t *testing.T) {
	assert.Len(t, token.New(10), 10)
	assert.Len(t, token.New(0), 0)
	assert.Panics(t, func() {
		token.New(-1)
	})
}
