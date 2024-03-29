// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package webo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
)

type SessionStore struct {
	*sessions.Session
}

func (f *SessionStore) PutString(key string, value string) {
	f.Values[key] = value
}
func (f *SessionStore) GetString(key string) string {
	r, _ := f.Values[key].(string)
	return r
}
func (f *SessionStore) PutInt(key string, value int) {
	f.Values[key] = value
}
func (f *SessionStore) GetInt(key string) (int, error) {
	r, e := f.Values[key].(int)
	if !e {
		return r, fmt.Errorf("invalid type assertion %v not int %v", f.Values[key], e)
	}
	return r, nil
}
func (f *SessionStore) PutDate(key string, value time.Time) {
	f.Values[key] = value.Unix()
}
func (f *SessionStore) GetDate(key string) (time.Time, error) {
	r, e := f.Values[key].(int64)
	if !e {
		return time.Now(), fmt.Errorf("invalid type assertion %v not unix time.time %v", f.Values[key], e)
	}

	return time.Unix(r, 0), nil
}

func (f *SessionStore) PutForm(key string, u url.Values) {
	f.PutString(key, u.Encode())
}

func (f *SessionStore) GetForm(key string) (url.Values, error) {
	return url.ParseQuery(f.GetString(key))
}

func (f *SessionStore) AddFlashf(flag string, msg string, a ...interface{}) {
	f.AddFlash(flag, fmt.Sprintf(msg, a...))
}
func (f *SessionStore) AddFlash(flag string, msg string) {
	f.Session.AddFlash(msg, flag)
}
func (f *SessionStore) Flashes(flag string) []string {
	fls := []string{}
	for _, s := range f.Session.Flashes(flag) {
		st, _ := s.(string)
		fls = append(fls, st)
	}
	return fls
}
func (f *SessionStore) Info(msg string) { f.AddFlash("info", msg) }
func (f *SessionStore) Infof(msg string, a ...interface{}) {
	f.Info(fmt.Sprintf(msg, a...))
}
func (f *SessionStore) Infos() []string    { return f.Flashes("info") }
func (f *SessionStore) Warning(msg string) { f.AddFlash("warning", msg) }
func (f *SessionStore) Warningf(msg string, a ...interface{}) {
	f.Warning(fmt.Sprintf(msg, a...))
}
func (f *SessionStore) Warnings() []string { return f.Flashes("warning") }
func (f *SessionStore) Alert(msg string) {
	f.AddFlash("alert", msg)
}
func (f *SessionStore) Alertf(msg string, a ...interface{}) {
	f.Alert(fmt.Sprintf(msg, a...))
}
func (f *SessionStore) Alerts() []string { return f.Flashes("alert") }

func RequestSession(r *http.Request) *SessionStore {
	return r.Context().Value("webo-gsession").(*SessionStore)
}

type Session struct {
	next http.Handler
}

func (s *Session) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.next.ServeHTTP(w, r)
}

func NewSession(auth string, key string, age int, next http.Handler) *Session {
	return &Session{SessionMiddleware(auth, key, age)(next)}
}

// deprecate : use SessionNameMiddleware
func SessionMiddleware(auth string, key string, age int) func(next http.Handler) http.Handler {
	return SessionNameMiddleware("gos", auth, key, age)
}

// session with possibility to set the name
func SessionNameMiddleware(name string, auth string, key string, age int) func(next http.Handler) http.Handler {
	authb := make([]byte, 32, 32)
	keyb := make([]byte, 32, 32)
	copy(authb, []byte(auth))
	copy(keyb, []byte(key))
	store := sessions.NewCookieStore(authb, keyb)
	store.MaxAge(age)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, _ := store.Get(r, name)
			sesG := &SessionStore{s}
			ctx := context.WithValue(r.Context(), "webo-gsession", sesG)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			err := sesG.Save(r, w)
			if err != nil {
				log.Printf("session error : %s", err)
			}
		})
	}
}
