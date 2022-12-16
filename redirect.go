package webo

import "fmt"

type ErrRedirect struct {
	URL string
}

func (e ErrRedirect) Error() string {
	return e.URL
}

func Redirect(url string, args ...any) ErrRedirect {
	return ErrRedirect{fmt.Sprintf(url, args...)}
}
