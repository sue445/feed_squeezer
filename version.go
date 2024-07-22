package main

import (
	"fmt"
	"net/http"
)

// VersionHandler is handler for GET /api/version
func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, GetVersion())
}
