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

	renderTemplate(w, "public/index.html", data)
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	renderFile(w, "public/favicon.svg")
}

func renderFile(w http.ResponseWriter, filename string) {
	b, err := static.ReadFile(filename)
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderFile %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	content := string(b)
	fmt.Fprint(w, content)
}

func renderTemplate(w http.ResponseWriter, filename string, data any) {
	b, err := static.ReadFile(filename)
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderTemplate %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	tmpl := string(b)
	t, err := template.New(filename).Parse(tmpl)

	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderTemplate %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] renderTemplate %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, buf.String())
}
