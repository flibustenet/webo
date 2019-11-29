package webo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime/debug"

	"go.flib.fr/werr"
)

// RECOVER

// renvoi le log lui-même créé par catcher et positionné dans context
func RequestCatcherLog(r *http.Request) *log.Logger {
	return r.Context().Value("webo-catcher-log").(*log.Logger)
}

type Catcher struct {
	debug    int
	name     string
	url_log  string //envoi de l'erreur
	version  string
	poste    string
	send_log string // url pour send logs
	next     http.Handler
}

func (h *Catcher) ServeHTTP(wrt http.ResponseWriter, r *http.Request) {
	CatcherMiddleware(h.debug, h.name, h.url_log, h.version, h.poste)(h.next).ServeHTTP(wrt, r)
}

func NewCatcher(debug int, name string, url_log string, version string, poste string, h http.Handler) *Catcher {
	if version == "" {
		version = "0"
	}
	if poste == "" {
		poste = "no_poste"
	}
	if name == "" {
		name = "no_name"
	}
	//fmt.Printf("Start catcher [%s] debug:%d version:%s poste:%s url_log:%s", name, debug, version, poste, url_log)
	c := &Catcher{debug, name, url_log, version, poste, "", h}
	log.Printf("Start [%s] version:%s poste:%s\n", name, version, poste)
	fmt.Printf("Start [%s] version:%s poste:%s\n", name, version, poste)
	return c
}
func CatcherMiddleware(debugFlag int, name string, url_log string, version string, poste string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wrt http.ResponseWriter, r *http.Request) {
			// log pour la durée de la requête
			// affiché uniquement si erreur ou si debug=2
			var logBuf bytes.Buffer
			lg := log.New(&logBuf, "", log.LstdFlags)
			lg.Println("------", r.Method, r.URL)
			fmt.Println(r.Method, r.URL)
			ctx := context.WithValue(r.Context(), "webo-catcher-log", lg)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			defer func() {
				if rec := recover(); rec != nil {
					if rec_redirect, ok := rec.(ErrRedirect); ok {
						http.Redirect(wrt, r, rec_redirect.URL, http.StatusSeeOther)
						return
					}
					lg.Println("=== Panic ===")
					sdebug := ""
					switch x := rec.(type) {
					case error:
						sdebug = werr.SprintSkip(x, "ServeHTTP")
					default:
						sdebug = string(debug.Stack())
					}
					lg.Println(sdebug)

					log.Println(logBuf.String())
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
						fmt.Fprintln(wrt, fmt.Sprintf("<html><pre><b>%v</b></pre>", rec)+"<pre>"+logBuf.String()+"</pre>")
					} else {
						http.Error(wrt, "Travaux en cours...", http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(w, r)
			for k, v := range w.Header() {
				wrt.Header()[k] = v
			}
			if wrt.Header().Get("Content-Type") == "" {
				wrt.Header().Set("Content-Type", "text/html; charset=utf-8")
			}
			if w.Code == 0 {
				w.Code = 200
			}
			wrt.WriteHeader(w.Code)
			wrt.Write(w.Body.Bytes())
			// ignore broken pipe (client déjà fermé ?)
			if debugFlag == 2 {
				log.Print(logBuf.String())
			}
		})
	}
}
