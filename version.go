package main

import (
	"fmt"
	"net/http"
)

func versionHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, getVersion())
}
