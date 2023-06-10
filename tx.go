// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package webo

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type TX struct {
	DB   *sqlx.DB
	next http.Handler
}

func RequestTx(r *http.Request) *sqlx.Tx {
	return r.Context().Value("webo-tx").(*sqlx.Tx)
}

func (t *TX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	TxMiddleware(t.DB)(t.next).ServeHTTP(w, r)
}
func NewTx(db *sqlx.DB, h http.Handler) *TX {
	c := &TX{db, h}
	return c
}

func TxMiddleware(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := RequestCatcherLog(r)
			tx := db.MustBegin()

			defer func() {
				if rec := recover(); rec != nil {
					log.Print("rollback")
					tx.Rollback()
					panic(rec)
				} else {
					tx.Commit()
					log.Println("commit")
				}
			}()

			ctx := context.WithValue(r.Context(), "webo-tx", tx)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
