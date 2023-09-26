// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package webo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"runtime/debug"
)

// NOT TESTED
// DirectCatcherMiddleware
// sans httptest.NewRecord, envoi direct vers http.ResponseWriter
// possible http: superfluous response.WriteHeader à ignorer
func DirectCatcherMiddleware(debugFlag int, name string, url_log string, version string, poste string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wrt http.ResponseWriter, r *http.Request) {
			// log pour la durée de la requête
			// affiché uniquement si erreur ou si debug=2
			var logBuf bytes.Buffer
			lg := log.New(&logBuf, "", log.LstdFlags)
			lg.Println("------", r.Method, r.URL)
			lg.SetFlags(log.Lshortfile)
			//fmt.Println(r.Method, r.URL)
			ctx := context.WithValue(r.Context(), "webo-catcher-log", lg)
			r = r.WithContext(ctx)
			defer func() {
				if rec := recover(); rec != nil {
					if rec_redirect, ok := rec.(ErrRedirect); ok {
						http.Redirect(wrt, r, rec_redirect.URL, http.StatusSeeOther)
						return
					}
					lg.SetFlags(0)
					lg.Printf("panic: %s %s", r.Method, r.URL)
					sdebug := ""
					switch x := rec.(type) {
					case runtime.Error:
						sdebug = fmt.Sprintf("%s\n%v", debug.Stack(), rec)
					case error:
						sdebug = x.Error()
					default:
						sdebug = fmt.Sprintf("%s\n%v", debug.Stack(), rec)
					}
					lg.Println(sdebug)

					fmt.Println(logBuf.String())
					if debugFlag == 0 && url_log != "" {
						resp, err := http.PostForm(url_log, url.Values{"title": {"[bug] " + name + "_" + version},
							"version": {version},
							"poste":   {poste},
							"log":     {logBuf.String()}})

						if err != nil {
							log.Println(err)
							http.Error(wrt, "Travaux en cours ! ", http.StatusInternalServerError)
							return
						}
						http.Error(wrt, "Maintenance en cours... ", http.StatusInternalServerError)
						defer resp.Body.Close()
						return
					}
					if debugFlag > 0 {
						http.Error(wrt, logBuf.String(), http.StatusInternalServerError)
					} else {
						http.Error(wrt, "Travaux en cours...", http.StatusInternalServerError)
					}
				}
			}()
			wrt.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
			next.ServeHTTP(wrt, r)
		})
	}
}
