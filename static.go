package main

import (
	"embed"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"net/http"
)

//go:embed public/*
var static embed.FS

func indexHandler(res http.ResponseWriter, req *http.Request) {
	renderFile(res, "public/index.html")
}

func renderFile(res http.ResponseWriter, filename string) {
	b, err := static.ReadFile(filename)
	if err != nil {
		sentry.CaptureException(errors.WithStack(err))
		http.Error(res, "error", http.StatusInternalServerError)
	}

	content := string(b)
	fmt.Fprint(res, content)
	res.WriteHeader(http.StatusOK)
}
