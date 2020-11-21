package webo

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/sessions"
)

func TestSession(t *testing.T) {
	store := sessions.NewCookieStore([]byte("abcd"))
	r := httptest.NewRequest("GET", "/", nil)
	s, _ := store.Get(r, "gos")
	ses := &SessionStore{s}
	ses.PutInt("ok", 5)
	i, e := ses.GetInt("ok")
	if i != 5 {
		t.Errorf("ok=%d %v", i, e)
	}
	ses.Alert("ok")
	fls := ses.Flashes("alert")
	if len(fls) != 1 || fls[0] != "ok" {
		t.Errorf("alert=%s", ses.Flashes("alert"))
	}

	now := time.Now().Truncate(time.Second)
	ses.PutDate("d", now)
	d, _ := ses.GetDate("d")
	if !d.Equal(now) {
		t.Errorf("put date : %v %v", now, d)
	}
	now = time.Now()
	ses.PutDate("d", now)
	_, e = ses.GetInt("j")
	if e == nil {
		t.Errorf("should be nil")
	}

	u := url.Values{}
	u.Set("x", "xxx")
	u.Set("y", "yyy")
	ses.PutForm("fo", u)
	res, err := ses.GetForm("fo")
	if err != nil {
		t.Errorf("GetForm error : %v", err)
	}
	if res.Get("x") != "xxx" {
		t.Errorf("GetForm x should be xxx : %s", res["x"])
	}
	if res.Get("y") != "yyy" {
		t.Errorf("GetForm y should be yyy : %s", res["y"])
	}
}
