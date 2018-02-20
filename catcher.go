package webo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	"github.com/pkg/errors"
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
	var logBuf bytes.Buffer
	lg := log.New(&logBuf, "", log.LstdFlags)
	lg.Println("------", r.Method, r.URL)
	fmt.Println(r.Method, r.URL)
	if r.Method == "POST" {
		r.ParseForm()
		lg.Print("POST = ", r.Form)
	}
	ctx := context.WithValue(r.Context(), "webo-catcher-log", lg)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
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
			fmt.Println(logBuf.String())
			if h.debug == 0 && h.url_log != "" {
				resp, err := http.PostForm(h.url_log, url.Values{"title": {"[bug] " + h.name + "_" + h.version},
					"version": {h.version},
					"poste":   {h.poste},
					"log":     {logBuf.String()}})

				if err != nil {
					log.Println(err)
					http.Error(wrt, "500: Un incident s'est produit et n'a pas pu être envoyé à l'administrateur",
						http.StatusInternalServerError)
					return
				}
				http.Error(wrt, "500: Un incident s'est produit et a été envoyé à l'administrateur",
					http.StatusInternalServerError)
				defer resp.Body.Close()
				return
			}
			if h.debug > 0 {
				fmt.Fprintln(wrt, "<pre>"+logBuf.String()+"</pre>")
			} else {
				http.Error(wrt, "500: Un incident s'est produit",
					http.StatusInternalServerError)
			}
		}
	}()

	h.next.ServeHTTP(w, r)
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
	log.Print(logBuf.String())
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
