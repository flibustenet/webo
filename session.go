package webo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"
)

type SessionStore struct {
	sync.Mutex
	kv     map[string]string
	secure *securecookie.SecureCookie
}

func (f *SessionStore) PutString(key string, value string) {
	f.Lock()
	defer f.Unlock()
	f.kv["gos-"+key] = value
}
func (f *SessionStore) GetString(key string) string {
	f.Lock()
	defer f.Unlock()
	return f.kv["gos-"+key]
}
func (f *SessionStore) PutInt(key string, value int) {
	f.Lock()
	defer f.Unlock()
	f.kv["gos-"+key] = strconv.Itoa(value)
}
func (f *SessionStore) GetInt(key string) (int, error) {
	f.Lock()
	defer f.Unlock()
	res, ok := f.kv["gos-"+key]
	if !ok {
		res = "0"
	}
	resi, err := strconv.Atoi(res)
	if err != nil {
		return 0, errors.Wrapf(err, "cookie non int %s : %s", key, res)
	}
	return resi, err
}
func (f *SessionStore) PopInt(key string) (int, error) {
	resi, err := f.GetInt(key)
	if err != nil {
		return 0, err
	}
	f.Lock()
	defer f.Unlock()
	f.kv["gos-"+key] = ""
	return resi, nil
}

func (f *SessionStore) AddFlash(flag string, msg string) {
	f.Lock()
	defer f.Unlock()
	if msg == "" {
		return
	}
	kv := f.kv["gos-flash-"+flag]
	if kv != "" {
		kv += "§"
	}
	kv += msg
	f.kv["gos-flash-"+flag] = kv
}
func (f *SessionStore) Flashes(flag string) []string {
	f.Lock()
	defer f.Unlock()
	kv := f.kv["gos-flash-"+flag]
	if kv == "" {
		return make([]string, 0)
	}
	r := strings.Split(kv, "§")
	f.kv["gos-flash-"+flag] = ""
	return r
}
func (f *SessionStore) Info(msg string)    { f.AddFlash("info", msg) }
func (f *SessionStore) Infos() []string    { return f.Flashes("info") }
func (f *SessionStore) Warning(msg string) { f.AddFlash("warning", msg) }
func (f *SessionStore) Warnings() []string { return f.Flashes("warning") }
func (f *SessionStore) Alert(msg string)   { f.AddFlash("alert", msg) }
func (f *SessionStore) Alerts() []string   { return f.Flashes("alert") }

func (f *SessionStore) SetCookies(w http.ResponseWriter) error {
	log.Printf("set cookies %v", f.kv)
	for k, v := range f.kv {
		encode, err := f.secure.Encode(k, v)
		if err != nil {
			return errors.Wrapf(err, "Impossible d'enregistrer le cookie %s : %s", k, v)
		}
		log.Println("set ", k, v)
		http.SetCookie(w, &http.Cookie{Name: k, Value: encode, Path: "/"})
	}
	return nil
}
func RequestSession(r *http.Request) *SessionStore {
	return r.Context().Value("webo-session").(*SessionStore)
}

type Session struct {
	next http.Handler
}

func (s *Session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.next.ServeHTTP(w, r)
}

func NewSession(hmac string, bloc string, next http.Handler) *Session {
	return &Session{SessionMiddleware(hmac, bloc)(next)}
}

func SessionMiddleware(hmac string, bloc string) func(next http.Handler) http.Handler {
	hmacKey := []byte(fmt.Sprintf("%32s", hmac))
	blockKey := []byte(fmt.Sprintf("%32s", bloc))
	secure := securecookie.New(hmacKey, blockKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ses := &SessionStore{}
			ses.secure = secure
			ses.kv = make(map[string]string)
			for _, c := range r.Cookies() {
				if strings.HasPrefix(c.Name, "gos-") {
					var value string
					err := secure.Decode(c.Name, c.Value, &value)
					if err != nil {
						log.Printf("Impossible de décoder le cookie %s : %s", c.Name, c.Value)
					}
					ses.Lock()
					ses.kv[c.Name] = value
					ses.Unlock()
				}
			}
			ctx := context.WithValue(r.Context(), "webo-session", ses)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			ses.SetCookies(w)
		})
	}
}
