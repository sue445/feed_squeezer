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
	w := httptest.NewRecorder()

	main.IndexHandler(w, r)

	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if assert.NoError(t, err) {
		assert.Contains(t, string(body), "<title>feed_squeezer</title>")
	}
}

func TestIndexHandler_NotFound(t *testing.T) {
	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()

	main.IndexHandler(w, r)

	res := w.Result()

	assert.Equal(t, 404, res.StatusCode)
}
