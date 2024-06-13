package main

import (
	"fmt"
	"net/http"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, getVersion())
}
