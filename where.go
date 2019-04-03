package webo

import (
	"fmt"
	"strings"
)

type Where struct {
	Style string // $ ou ?
	where []string
	Args  []interface{}
}

//
// remplace %s par ? ou $d
//
func (w *Where) Appendf(s string, a ...interface{}) {
	if len(s) > 0 {
		arg_nb := []interface{}{} // deviendra $1 $2 $...
		for i := 0; i < len(a); i++ {
			switch w.Style {
			case "?":
				arg_nb = append(arg_nb, "?")
			case "", "$":
				arg_nb = append(arg_nb, fmt.Sprintf("$%d", len(w.Args)+1+i))
			}
		}
		w.where = append(w.where, fmt.Sprintf(s, arg_nb...))
	}

	if len(a) > 0 {
		w.Args = append(w.Args, a...)
	}
}

//
// ajout sous forme de liste
// exemple AppendListf("xyz in (%s)", "a","b","c")
// doit ajouter "xyz in ($1,$2,$3)" avec args "a","b","c"
//
func (w *Where) AppendListf(s string, a ...interface{}) {
	q := []string{} // les $1 $2...
	for i := 0; i < len(a); i++ {
		if w.Style == "?" {
			q = append(q, "?")
			continue
		}
		w.Args = append(w.Args, a[i])
		q = append(q, fmt.Sprintf("$%d", len(w.Args)))
	}
	w.where = append(w.where, fmt.Sprintf(s, strings.Join(q, ",")))
}

//
// renvoi suite avec and, sans le premier and
//
func (w *Where) And() string {
	return strings.Join(w.where, " and ")
}

//
// renvoi suite de and avec le premier and si non vide
//
func (w *Where) AndAnd() string {
	if len(w.where) == 0 {
		return ""
	}
	return " and " + w.And()
}

//
// renvoi le where avec "where" sauf si vide
//
func (w *Where) Where() string {
	if len(w.where) == 0 {
		return ""
	}
	ands := w.And()
	return " where " + ands
}

//
// renvoi une copie
//
func (w *Where) Clone() *Where {
	wc := &Where{}
	wc.Args = append(wc.Args, w.Args...)
	wc.where = append(wc.where, w.where...)
	return wc
}
