package main_test

import (
	"flag"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/feed_squeezer"
	"os"
	"testing"
)

func init() {
	testing.Init()
	flag.Parse()
}

func TestGenerateSqueezedAtom(t *testing.T) {
	type args struct {
		feedXMLFile string
		query       string
	}
	tests := []struct {
		name        string
		args        args
		wantXMLFile string
	}{
		{
			name: "Normal",
			args: args{
				feedXMLFile: "testdata/youtube_toei_animation.atom",
				query:       "おジャ魔女どれみ",
			},
			wantXMLFile: "testdata/youtube_toei_animation_ojamajo.atom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed := ReadTestData(tt.args.feedXMLFile)
			want := ReadTestData(tt.wantXMLFile)

			got, err := main.GenerateSqueezedAtom(feed, tt.args.query)

			if assert.NoError(t, err) {
				assert.Equal(t, want, got)
			}
		})
	}
}

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

func TestGetContentFromCache(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.com/test.txt",
		httpmock.NewStringResponder(200, ReadTestData("testdata/test.txt")))

	gotFromOrigin, err := main.GetContentFromURL("http://example.com/test.txt")
	if assert.NoError(t, err) {
		assert.Equal(t, gotFromOrigin, "test")
	}

	gotFromCache, err := main.GetContentFromURL("http://example.com/test.txt")
	if assert.NoError(t, err) {
		assert.Equal(t, gotFromCache, "test")
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

func TestNormalize(t *testing.T) {
	tests := []struct {
		str  string
		want string
	}{
		{
			str:  "GitLab",
			want: "gitlab",
		},
		{
			str:  "１２３",
			want: "123",
		},
		{
			str:  " abc ",
			want: "abc",
		},
		{
			str:  " ",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := main.Normalize(tt.str)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		message string
		want    int
	}{
		{
			message: "404 Not Found",
			want:    404,
		},
		{
			message: "Not Found",
			want:    -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got, err := main.GetStatusCode(tt.message)

			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
