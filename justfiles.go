package webo

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func StaticMiddleware(r *mux.Router, static string, maxAge int) func(next http.Handler) http.Handler {
	hdl := http.StripPrefix("/"+static, http.FileServer(JustFiles{http.Dir(static)}))
	max_age := strconv.Itoa(maxAge)
	r.PathPrefix("/" + static).Handler(hdl)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.RequestURI, "/"+static) {
				w.Header().Set("Cache-Control", "max-age="+max_age)
				hdl.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func StaticFsMiddleware(r *mux.Router, static string, fs http.FileSystem, maxAge int) func(next http.Handler) http.Handler {
	hdl := http.FileServer(fs)
	max_age := strconv.Itoa(maxAge)
	r.PathPrefix("/" + static).Handler(hdl)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.RequestURI, "/"+static) {
				w.Header().Set("Cache-Control", "max-age="+max_age)
				hdl.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Assets(r *mux.Router, static string) {
	r.PathPrefix("/" + static).Handler(http.StripPrefix("/"+static, http.FileServer(JustFiles{http.Dir(static)}))).Name(static)
	//r.PathPrefix("/" + static).Handler(http.StripPrefix("/"+static, http.FileServer(http.Dir(static)))).Name(static)
}

func AssetsDir(r *mux.Router, static string, dir string) {
	r.PathPrefix(static).Handler(http.StripPrefix(static, http.FileServer(JustFiles{http.Dir(dir)}))).Name(static)
	//r.PathPrefix("/" + static).Handler(http.StripPrefix("/"+static, http.FileServer(http.Dir(static)))).Name(static)
}

type JustFiles struct {
	FS http.FileSystem
}

func (fs JustFiles) Open(name string) (http.File, error) {
	f, err := fs.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
