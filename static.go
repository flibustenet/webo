package webo

import (
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type static struct {
	path       string
	maxAge     int
	origin     string
	fileServer http.Handler // http.FileServer
}

func (s static) handler(r *mux.Router) func(next http.Handler) http.Handler {
	path := s.path
	if path == "" {
		path = "/static"
	}
	if len(path) > 0 && path[0] != '/' {
		path = "/" + s.path
	}
	hdl := http.StripPrefix(path, s.fileServer)
	max_age := strconv.Itoa(s.maxAge)
	// must inform router that this path exists
	r.PathPrefix(path).Handler(hdl)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.RequestURI, path+"/") {
				//	log.Println(r.Method, r.RequestURI)
				w.Header().Set("Cache-Control", "max-age="+max_age)
				if s.origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", s.origin)
				}
				hdl.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// with gorilla mux.Router (gorilla must know that /static exists)
// add a static middleware from disk
// with :
// Cache-Control max-age=maxAge
func StaticDiskMiddleware(r *mux.Router, path string, maxAge int) func(next http.Handler) http.Handler {
	return StaticDiskMiddlewareOrigin(r, path, maxAge, "")
}

// with gorilla mux.Router (gorilla must know that /static exists)
// add a static middleware from disk
// with :
// Cache-Control max-age=maxAge
// Access-Control-Allow-Origin origin
func StaticDiskMiddlewareOrigin(r *mux.Router, path string, maxAge int, origin string) func(next http.Handler) http.Handler {
	return static{
		path:       path,
		maxAge:     maxAge,
		origin:     origin,
		fileServer: http.FileServer(JustFiles{http.Dir(path)}),
	}.handler(r)
}

// with gorilla mux.Router (gorilla must know that /static exists)
// add a static middleware from fs.FS
// with :
// Cache-Control maxAge
func StaticFsMiddleware(r *mux.Router, fs fs.FS, path string, maxAge int) func(next http.Handler) http.Handler {
	return StaticFsMiddlewareOrigin(r, fs, path, maxAge, "")
}

// with gorilla mux.Router (gorilla must know that /static exists)
// add a static middleware from fs.FS
// with :
// Cache-Control max-age=maxAge
// Access-Control-Allow-Origin origin
func StaticFsMiddlewareOrigin(r *mux.Router, fs fs.FS, path string, maxAge int, origin string) func(next http.Handler) http.Handler {
	return static{
		path:       path,
		maxAge:     maxAge,
		origin:     origin,
		fileServer: http.FileServer(http.FS(fs)),
	}.handler(r)
}
