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
	renderFile(w, "public/index.html")
}

func renderFile(w http.ResponseWriter, filename string) {
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
	}{
		Version: GetVersion(),
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
