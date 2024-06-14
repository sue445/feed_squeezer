package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"text/template"
)

//go:embed public/*
var static embed.FS

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderFile(w, r, "public/index.html")
}

func renderFile(w http.ResponseWriter, r *http.Request, filename string) {
	b, err := static.ReadFile(filename)
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderFile %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	tmpl := string(b)
	t, err := template.New("index").Parse(tmpl)

	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderFile %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Version string
		BaseURL string
	}{
		Version: GetVersion(),
		BaseURL: getBaseURL(r),
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderFile %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, buf.String())
}

func getBaseURL(r *http.Request) string {
	return fmt.Sprintf("%s://%s/", getScheme(r), r.Host)
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}

	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}

	return "http"
}
