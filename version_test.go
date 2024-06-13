package main_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/feed_squeezer"
	"io"
	"net/http/httptest"
	"testing"
)

func TestVersionHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/api/version", nil)
	w := httptest.NewRecorder()

	main.VersionHandler(w, r)

	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if assert.NoError(t, err) {
		assert.Contains(t, string(body), main.GetVersion())
	}
}
