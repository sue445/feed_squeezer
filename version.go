package main

import (
	"fmt"
	"net/http"
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, getVersion())
}
