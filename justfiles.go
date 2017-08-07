package webo

import (
	"net/http"
	"os"
)

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
