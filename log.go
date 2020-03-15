package webo

import (
	"log"
	"net/http"
)

func Log(r *http.Request) {
	log.Println(r.Method, r.URL)
}
