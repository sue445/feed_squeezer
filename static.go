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

	data := struct {
		Version string
	}{
		Version: GetVersion(),
	}

	content, err := renderTemplate("public/index.html", data)
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] IndexHandler %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, content)
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	content, err := renderFile("public/favicon.svg")
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] FaviconHandler %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// NOTE: Response headers must be set before writing to response body
	w.Header().Set("Content-Type", "image/svg+xml")

	fmt.Fprint(w, content)
}

func renderFile(filename string) (string, error) {
	b, err := static.ReadFile(filename)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(b), nil
}

func renderTemplate(filename string, data any) (string, error) {
	b, err := static.ReadFile(filename)
	if err != nil {
		return "", errors.WithStack(err)
	}

	tmpl := string(b)
	t, err := template.New(filename).Parse(tmpl)
	if err != nil {
		return "", errors.WithStack(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return buf.String(), nil
}
