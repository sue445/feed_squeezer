package main_test

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/feed_proxy"
	"os"
	"testing"
)

func TestGetContentFromURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.com/test.txt",
		httpmock.NewStringResponder(200, ReadTestData("testdata/test.txt")))

	got, err := main.GetContentFromURL("http://example.com/test.txt")
	if assert.NoError(t, err) {
		assert.Equal(t, got, "test")
	}
}

// ReadTestData returns testdata
func ReadTestData(filename string) string {
	buf, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(buf)
}
