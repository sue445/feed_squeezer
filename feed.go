package main

import (
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"io"
	"net/http"
	"strconv"
)

// GetContentFromURL returns content from URL (without cache)
func GetContentFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.WithStack(err)
	}

	defer resp.Body.Close()

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(map[string]string{
			"Status":     resp.Status,
			"StatusCode": strconv.Itoa(resp.StatusCode),
		})
	})

	if resp.StatusCode >= 400 {
		return "", errors.New(resp.Status)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(buf), nil
}
