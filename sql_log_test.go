package webo

import (
	"log"
	"testing"
	"time"
)

func Test_sql_quote(t *testing.T) {
	type Tst struct {
		T int
		V interface{}
		S string
	}
	tbl := []Tst{Tst{DB_ACCESS, 5, "5"},
		Tst{DB_ACCESS, 3.1415, "3.1415"},
		Tst{DB_ACCESS, 3., "3"},
		Tst{DB_ACCESS, "abcd", "'abcd'"},
		Tst{DB_ACCESS, "ab'cd", "'ab''cd'"},
		Tst{DB_ACCESS, true, "-1"},
		Tst{DB_ACCESS, false, "0"},
		Tst{DB_ACCESS, nil, "null"},
		Tst{DB_ACCESS, time.Date(1969, 11, 05, 23, 05, 03, 0, time.Local), "'1969-11-05 23:05:03'"},
		Tst{DB_PG, time.Date(1969, 11, 05, 23, 05, 03, 0, time.Local), "'1969-11-05 23:05:03'"},
	}
	for _, s := range tbl {
		r := sql_quoter(s.T, s.V)
		if r != s.S {
			log.Fatalf("attend %s reçoit %s", s.S, r)
		}
	}
}
func Test_sql_quote_query(t *testing.T) {
	type Tst struct {
		T int
		Q string
		V []interface{}
		S string
	}
	tbl := []Tst{
		Tst{DB_ACCESS, "? ? ? ? ? ?", []interface{}{5, "abcd", "e'fg", true, false, nil}, "5 'abcd' 'e''fg' -1 0 null"},
		Tst{DB_PG, "$1, $2, $3 $4 $5 $6", []interface{}{5, "abcd", "e'fg", true, false, nil}, "5, 'abcd', 'e''fg' true false null"},
		Tst{DB_PG, "$1 $3 $2 $3", []interface{}{5, "abcd", "e'fg"}, "5 'e''fg' 'abcd' 'e''fg'"},
	}
	for _, s := range tbl {
		r := sql_fake(s.T, s.Q, s.V...)
		if r != s.S {
			log.Fatalf("type %d : %s attend %s reçoit %s", s.T, s.Q, s.S, r)
		}
	}
}
