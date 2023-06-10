// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package webo

import (
	"log"
	"net/http"
)

func Log(r *http.Request) {
	log.Println(r.Method, r.URL)
}

// HTTP middleware setting a value on the request context
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log(r)
		next.ServeHTTP(w, r)
	})
}
