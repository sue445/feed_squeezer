package main_test

import (
	"flag"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/feed_squeezer"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		{
			name: "keyword included in content tag",
			args: args{
				feedXMLFile: "testdata/functions-release-notes.xml",
				query:       "ruby",
			},
			wantXMLFile: "testdata/functions-release-notes-ruby.xml",
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
	httpmock.Activate(t)

	httpmock.RegisterResponder("GET", "http://example.com/test.txt",
		httpmock.NewStringResponder(200, ReadTestData("testdata/test.txt")))

	got, err := main.GetContentFromURL("http://example.com/test.txt")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", got)
	}
}

func TestGetContentFromURLGoogleCloudFeed(t *testing.T) {
	httpmock.Activate(t)
	originalVersion := main.Version
	main.Version = "v1.2.3"
	t.Cleanup(func() { main.Version = originalVersion })

	const feedURL = "https://cloud.google.com/feeds/functions-release-notes.xml"
	const redirectedURL = "https://docs.cloud.google.com/feeds/functions-release-notes.xml"
	want := ReadTestData("testdata/functions-release-notes.xml")

	checkUserAgent := func(t *testing.T, r *http.Request) bool {
		return assert.Equal(t, "feed_squeezer/v1.2.3 (+https://github.com/sue445/feed_squeezer)", r.Header.Get("User-Agent"))
	}
	httpmock.RegisterResponder("GET", feedURL, func(r *http.Request) (*http.Response, error) {
		if !checkUserAgent(t, r) {
			return httpmock.NewStringResponse(http.StatusOK, "OK"), nil
		}
		response := httpmock.NewStringResponse(http.StatusMovedPermanently, "")
		response.Header.Set("Location", redirectedURL)
		return response, nil
	})
	httpmock.RegisterResponder("GET", redirectedURL, func(r *http.Request) (*http.Response, error) {
		checkUserAgent(t, r)
		return httpmock.NewStringResponse(http.StatusOK, want), nil
	})

	got, err := main.GetContentFromURL(feedURL)
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
		atom, err := main.GenerateSqueezedAtom(got, "ruby")
		if assert.NoError(t, err) {
			assert.Equal(t, ReadTestData("testdata/functions-release-notes-ruby.xml"), atom)
		}
	}
}

func TestGetContentFromCache(t *testing.T) {
	httpmock.Activate(t)

	httpmock.RegisterResponder("GET", "http://example.com/test.txt",
		httpmock.NewStringResponder(200, ReadTestData("testdata/test.txt")))

	gotFromOrigin, err := main.GetContentFromURL("http://example.com/test.txt")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", gotFromOrigin)
	}

	gotFromCache, err := main.GetContentFromURL("http://example.com/test.txt")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", gotFromCache)
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

func TestFeedHandler(t *testing.T) {
	httpmock.Activate(t)

	httpmock.RegisterResponder("GET", "http://example.com/feed",
		httpmock.NewStringResponder(200, ReadTestData("testdata/youtube_toei_animation.atom")))

	r := httptest.NewRequest("GET", "/api/feed", nil)
	w := httptest.NewRecorder()

	r.Form = url.Values{}
	r.Form.Add("url", "http://example.com/feed")
	r.Form.Add("query", "おジャ魔女どれみ")

	main.FeedHandler(w, r)

	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if assert.NoError(t, err) {
		assert.Contains(t, string(body), ReadTestData("testdata/youtube_toei_animation_ojamajo.atom"))
	}
}
