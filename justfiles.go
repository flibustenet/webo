package webo

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

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
