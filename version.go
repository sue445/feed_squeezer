package main

import (
	"fmt"
	"net/http"
)

func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, GetVersion())
}
