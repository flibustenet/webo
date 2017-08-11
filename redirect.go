package webo

type ErrRedirect struct {
	URL string
}

func (e ErrRedirect) Error() string {
	return e.URL
}

func Redirect(url string) ErrRedirect {
	return ErrRedirect{url}
}
