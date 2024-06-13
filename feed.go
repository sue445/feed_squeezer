package main

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"golang.org/x/text/width"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

func feedHandler(w http.ResponseWriter, r *http.Request) {
	feedURL := r.FormValue("url")
	query := r.FormValue("query")

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(map[string]string{
			"feed":  feedURL,
			"query": query,
		})
	})

	atom, err := generateFeed(feedURL, query)

	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] feedHandler %v\n", errors.WithStack(err))

		status, err2 := GetStatusCode(err.Error())
		if err2 == nil && status >= 0 {
			// respect status code in error
			http.Error(w, err.Error(), status)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	fmt.Fprint(w, atom)
}

// GetStatusCode GetStatusCode returns status code in message
func GetStatusCode(message string) (int, error) {
	re := regexp.MustCompile(`^\d{3} `)
	match := re.FindString(message)

	if match == "" {
		return -1, nil
	}

	code, err := strconv.Atoi(strings.TrimSpace(match))
	if err != nil {
		return -1, err
	}
	return code, nil
}

func generateFeed(feedURL string, query string) (string, error) {
	feedData, err := GetContentFromCache(feedURL)
	if err != nil {
		return "", errors.WithStack(err)
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("generateFeed", map[string]interface{}{
			"feedData": feedData,
		})
	})

	atom, err := GenerateSqueezedAtom(feedData, query)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return atom, nil
}

// GenerateSqueezedAtom squeeze feedData with query
func GenerateSqueezedAtom(feedData string, query string) (string, error) {
	fp := gofeed.NewParser()

	feed, err := fp.Parse(strings.NewReader(feedData))
	if err != nil {
		return "", errors.WithStack(err)
	}

	result := feeds.Feed{
		Title:       fmt.Sprintf("%s (query:%s)", feed.Title, query),
		Link:        &feeds.Link{Href: feed.Link},
		Description: feed.Description,
		Items:       nil,
		Copyright:   feed.Copyright,
	}

	var itemUpdatedTimes []*time.Time

	query = Normalize(query)

	for _, item := range feed.Items {
		itemUpdatedTimes = append(itemUpdatedTimes, item.UpdatedParsed)
		description := getItemDescription(item)

		text := item.Title + " " + description

		contained, err := ContainsKeyword(Normalize(text), query)
		if err != nil {
			return "", errors.WithStack(err)
		}

		if contained {
			result.Items = append(result.Items, &feeds.Item{
				Id:          item.GUID,
				Title:       item.Title,
				Link:        &feeds.Link{Href: item.Link},
				Description: description,
				Created:     *item.PublishedParsed,
			})
		}
	}

	latestUpdatedTime := maxTime(itemUpdatedTimes)
	if latestUpdatedTime != nil {
		result.Updated = *latestUpdatedTime
	}

	atom, err := result.ToAtom()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return atom + "\n", nil
}

// Normalize normalize string
func Normalize(str string) string {
	str = width.Fold.String(str)
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	return str
}

func getItemDescription(item *gofeed.Item) string {
	// Find <media:description>
	groups := item.Extensions["media"]["group"]
	for _, group := range groups {
		if len(group.Children["description"]) > 0 {
			return group.Children["description"][0].Value
		}
	}

	return item.Description
}

func maxTime(times []*time.Time) *time.Time {
	if len(times) < 1 {
		return nil
	}

	max := times[0]
	for _, t := range times {
		if t.After(*max) {
			max = t
		}
	}
	return max
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
		log.Printf("[WARN] GetContentFromCache %v\n", errors.WithStack(err))

		// Suppress errors but send them to Sentry
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelWarning)

			sentry.CaptureException(errors.WithStack(err))
		})
	}

	// return from origin
	return content, nil
}
