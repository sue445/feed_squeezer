package main

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var feedCache *bigcache.BigCache

func init() {
	// NOTE: Slackbot calls GET immediately after HEAD for the same URL, so it is easy to get caught by YouTube's RateLimit, so cache the feed for a short time.
	config := bigcache.DefaultConfig(1 * time.Minute)

	var err error
	feedCache, err = bigcache.New(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
}

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

// GetContentFromCache returns content from URL (with cache)
func GetContentFromCache(url string) (string, error) {
	if entry, err := feedCache.Get(url); err == nil {
		// return from cache
		return string(entry), nil
	}

	content, err := GetContentFromURL(url)
	if err != nil {
		return "", errors.WithStack(err)
	}

	err = feedCache.Set(url, []byte(content))
	if err != nil {
		// Suppress errors but send them to Sentry
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelWarning)

			sentry.CaptureException(errors.WithStack(err))
		})
	}

	// return from origin
	return content, nil
}
