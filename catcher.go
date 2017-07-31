package webo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// RECOVER

// renvoi le log lui-même créé par catcher et positionné dans context
func CatcherLog(r *http.Request) *log.Logger {
	return r.Context().Value("webo-catcher-log").(*log.Logger)
}

type Catcher struct {
	debug   int
	name    string
	url_log string
	version string
	poste   string
	next    http.Handler
}

func (h *Catcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var logBuf bytes.Buffer
	lg := log.New(&logBuf, "", log.Lshortfile)
	lg.Println(r.Method, r.URL)
	ctx := context.WithValue(r.Context(), "webo-catcher-log", lg)
	r = r.WithContext(ctx)
	defer func() {
		if rec := recover(); rec != nil {
			lg.Println("================== Panic ==============")
			s := fmt.Sprintf("%+v", rec)
			if !strings.Contains(s, "runtime.") { // si pas de stack, créée une nouvelle erreur
				switch x := rec.(type) {
				case error:
					lg.Printf("%+v", errors.WithStack(rec.(error)))
				case string:
					lg.Printf("%+v", errors.New(x))
				default:
					lg.Printf("Unknow error type %s : %v", reflect.TypeOf(rec), rec)
				}
			} else {
				lg.Print(s)
			}
			log.Println(logBuf.String())
			if h.debug == 0 && h.url_log != "" {
				resp, err := http.PostForm(h.url_log, url.Values{"title": {"[bug] " + h.name},
					"version": {h.version},
					"poste":   {h.poste},
					"log":     {logBuf.String()}})

				if err != nil {
					log.Println(err)
					http.Error(w, "500: Un incident s'est produit et n'a pas pu être envoyé à l'administrateur",
						http.StatusInternalServerError)
					return
				}
				http.Error(w, "500: Un incident s'est produit et a été envoyé à l'administrateur",
					http.StatusInternalServerError)
				defer resp.Body.Close()
				return
			}
			if h.debug > 0 {
				fmt.Fprintf(w, logBuf.String())
			} else {
				http.Error(w, "500: Un incident s'est produit",
					http.StatusInternalServerError)
			}
		}
	}()

	h.next.ServeHTTP(w, r)
	fmt.Fprintf(os.Stderr, logBuf.String())
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
	c := &Catcher{debug, name, url_log, version, poste, h}
	fmt.Printf("Start catcher %#v", c) //[%s] debug:%d version:%s poste:%s url_log:%s", name, debug, version, poste, url_log)
	return c
}
