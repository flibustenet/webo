package webo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/engine/cookiestore"
	"github.com/alexedwards/scs/session"
)

type Session struct {
	Request *http.Request
}

func (f *Session) AddFlash(flag string, msg string) {
	if msg == "" {
		return
	}
	msgs, _ := session.GetString(f.Request, "webo-flashes-"+flag)

	if msgs != "" {
		msgs += "§"
	}
	msgs += msg
	session.PutString(f.Request, "webo-flashes-info", msgs)
}
func (f *Session) Flashes(flag string) []string {
	msgs, _ := session.GetString(f.Request, "webo-flashes-"+flag)
	r := strings.Split(msgs, "§")
	session.PutString(f.Request, "webo-flashes-info", "")
	return r
}
func (f *Session) Info(msg string)    { f.AddFlash("info", msg) }
func (f *Session) Infos() []string    { return f.Flashes("info") }
func (f *Session) Warning(msg string) { f.AddFlash("warning", msg) }
func (f *Session) Warnings() []string { return f.Flashes("warning") }
func (f *Session) Alert(msg string)   { f.AddFlash("alert", msg) }
func (f *Session) Alerts() []string   { return f.Flashes("alert") }

func RequestSession(r *http.Request) *Session {
	return r.Context().Value("webo-session").(*Session)
}

type SessionH struct {
	next http.Handler
}

func (s *SessionH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flash := &Session{Request: r}
	ctx := context.WithValue(r.Context(), "webo-flash", flash)
	r = r.WithContext(ctx)
	s.next.ServeHTTP(w, r)
}

func NewSessionCookieStore(hmac string, bloc string, next http.Handler) http.Handler {
	hmacKey := []byte(fmt.Sprintf("%32s", hmac))
	blockKey := []byte(fmt.Sprintf("%32s", bloc))
	keyset, err := cookiestore.NewKeyset(hmacKey, blockKey)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new CookieStore instance using the keyset.
	engine := cookiestore.New(keyset)
	sessionManager := session.Manage(engine)
	return sessionManager(&SessionH{next})
}
