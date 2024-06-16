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

	err := renderTemplate(w, "public/index.html", data)
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] IndexHandler %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	err := renderFile(w, "public/favicon.svg")
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		log.Printf("[ERROR] FaviconHandler %v\n", errors.WithStack(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
}

func renderFile(w http.ResponseWriter, filename string) error {
	b, err := static.ReadFile(filename)
	if err != nil {
		return err
	}

	content := string(b)
	fmt.Fprint(w, content)

	return nil
}

func renderTemplate(w http.ResponseWriter, filename string, data any) error {
	b, err := static.ReadFile(filename)
	if err != nil {
		return err
	}

	tmpl := string(b)
	t, err := template.New(filename).Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		return err
	}

	fmt.Fprint(w, buf.String())

	return nil
}
