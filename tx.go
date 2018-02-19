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
	log := RequestCatcherLog(r)
	tx := t.DB.MustBegin()

	defer func() {
		if rec := recover(); rec != nil {
			log.Print("rollback")
			tx.Rollback()
			panic(rec)
		} else {
			tx.Commit()
		}
	}()

	ctx := context.WithValue(r.Context(), "webo-tx", tx)
	r = r.WithContext(ctx)
	t.next.ServeHTTP(w, r)
}
func NewTx(db *sqlx.DB, h http.Handler) *TX {
	c := &TX{db, h}
	return c
}
