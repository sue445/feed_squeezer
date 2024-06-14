package main_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/feed_squeezer"
	"io"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler_Found(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Host = "localhost:8080"

	w := httptest.NewRecorder()

	main.IndexHandler(w, r)

	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	if assert.NoError(t, err) {
		body := string(b)
		assert.Contains(t, body, "<title>feed_squeezer</title>")
		assert.Contains(t, body, `<meta property="og:url" content="http://localhost:8080/" />`)
		assert.Contains(t, body, main.GetVersion())
	}
}

func TestIndexHandler_NotFound(t *testing.T) {
	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()

	main.IndexHandler(w, r)

	res := w.Result()

	assert.Equal(t, 404, res.StatusCode)
}
