package webo

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type Tx struct {
	*sqlx.Tx
	Mapper *reflectx.Mapper
}

func FromTx(tx *sqlx.Tx) *Tx {
	t := &Tx{Tx: tx}
	t.Mapper = tx.Mapper
	return t
}

func (x *Tx) InsertStruct(table string, f interface{}) (sql.Result, error) {
	return x.InsertMap(table, x.structToMap(f))
}
func (x *Tx) InsertStructExclude(table string, f interface{}, exclude []string) (sql.Result, error) {
	return x.InsertMap(table, x.structToMapExclude(f, exclude))
}
func (x *Tx) InsertStructInclude(table string, f interface{}, include []string) (sql.Result, error) {
	return x.InsertMap(table, x.structToMapInclude(f, include))
}

// renvoi la chaine sql et les valeurs pour un insert
// à partir d'un map
func (x *Tx) insertSt(table string, m map[string]interface{}) (string, []interface{}) {
	fieldols := make([]string, 0)
	values := make([]interface{}, 0)
	fieldnames := make([]string, 0)
	for name, _ := range m {
		fieldnames = append(fieldnames, name)
	}
	sort.Strings(fieldnames)

	for i, name := range fieldnames {
		fieldols = append(fieldols, fmt.Sprintf("$%d", i+1))
		values = append(values, m[name])
	}
	s := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fieldnames, ","),
		strings.Join(fieldols, ","))
	return s, values
}
func (x *Tx) InsertMap(table string, m map[string]interface{}) (sql.Result, error) {
	s, values := x.insertSt(table, m)
	res, err := x.Tx.Exec(s, values...)
	return res, err
}

func (x *Tx) UpdateStruct(table string, f interface{}, where string, where_vals ...interface{}) (sql.Result, error) {
	return x.UpdateMap(table, x.structToMap(f), where, where_vals...)
}
func (x *Tx) UpdateStructExclude(table string, f interface{}, exclude []string, where string, where_vals ...interface{}) (sql.Result, error) {
	return x.UpdateMap(table, x.structToMapExclude(f, exclude), where, where_vals...)
}
func (x *Tx) UpdateStructInclude(table string, f interface{}, include []string, where string, where_vals ...interface{}) (sql.Result, error) {
	return x.UpdateMap(table, x.structToMapInclude(f, include), where, where_vals...)
}

// renvoi la chaine sql et les valeurs pour un update
// à partir d'un map
func (x *Tx) updateSt(table string, m map[string]interface{}, where string, where_vals ...interface{}) (string, []interface{}) {
	sets := make([]string, 0)
	num := len(where_vals) + 1
	values := where_vals[:]

	fieldnames := make([]string, 0)
	for name, _ := range m {
		fieldnames = append(fieldnames, name)
	}
	sort.Strings(fieldnames)
	for _, name := range fieldnames {
		sets = append(sets, fmt.Sprintf("%s=$%d", name, num))
		num += 1
		values = append(values, m[name])
	}
	s := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(sets, ", "),
		where)
	return s, values
}
func (x *Tx) UpdateMap(table string, m map[string]interface{}, where string, where_vals ...interface{}) (sql.Result, error) {
	s, values := x.updateSt(table, m, where, where_vals...)
	log.Println(s)
	res, err := x.Tx.Exec(s, values...)

	return res, err
}

// convert struc to map
// avec mapper de sqlx.Tx
func (x *Tx) structToMap(f interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	fields := x.Mapper.FieldMap(reflect.ValueOf(f))
	for name, v := range fields {
		if strings.Contains(name, ".") {
			continue
		}
		m[name] = v.Interface()
	}
	return m
}
func (x *Tx) structToMapExclude(f interface{}, exclude []string) map[string]interface{} {
	m := x.structToMap(f)
	excludes := make(map[string]bool)
	for _, k := range exclude {
		excludes[k] = true
	}
	for k, _ := range m {
		if _, ok := excludes[k]; ok {
			delete(m, k)
		}
	}
	return m
}
func (x *Tx) structToMapInclude(f interface{}, include []string) map[string]interface{} {
	m := x.structToMap(f)
	includes := make(map[string]bool)
	for _, k := range include {
		includes[k] = true
	}
	for k, _ := range m {
		if _, ok := includes[k]; !ok {
			delete(m, k)
		}
	}
	return m
}
