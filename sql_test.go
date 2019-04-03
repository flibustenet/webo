package webo

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx/reflectx"
)

func Test_insertSt(t *testing.T) {
	tx := &Tx{}
	tx.Mapper = reflectx.NewMapperFunc("db", CamelToSnake)
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	q, _ := tx.insertSt("mytable", fs)
	if q != "INSERT INTO mytable (ok, yes) VALUES ($1, $2)" {
		log.Fatal(q)
	}
}

func Test_updateSt(t *testing.T) {
	tx := &Tx{}
	tx.Mapper = reflectx.NewMapperFunc("db", CamelToSnake)
	fs := map[string]interface{}{}
	fs["ok"] = "coral"
	fs["yes"] = "no"
	q, _ := tx.updateSt("mytable", fs, "ok=$1", "ok")
	if q != "UPDATE mytable SET ok=$2, yes=$3 WHERE ok=$1" {
		log.Fatal(q)
	}
}

func Test_structToMap(t *testing.T) {
	tx := &Tx{}
	tx.Mapper = reflectx.NewMapperFunc("db", CamelToSnake)
	type S struct {
		Ok     int
		Coral  string
		UnDeux string
	}
	m := tx.structToMap(&S{2, "coral", "un_deux"})
	if m["ok"].(int) != 2 {
		log.Fatal(m)
	}
	if m["coral"].(string) != "coral" {
		log.Fatal(m)
	}
	if m["un_deux"].(string) != "un_deux" {
		log.Fatal(m)
	}
}

func Test_structToMapInclude(t *testing.T) {
	tx := &Tx{}
	tx.Mapper = reflectx.NewMapperFunc("db", CamelToSnake)
	type S struct {
		Ok     int
		Coral  string
		UnDeux string
	}
	m := tx.structToMapInclude(&S{2, "coral", "un_deux"}, []string{"ok", "un_deux"})
	if m["ok"].(int) != 2 {
		log.Fatal(m)
	}
	if _, ok := m["coral"]; ok == true {
		log.Fatal(m)
	}
	if m["un_deux"].(string) != "un_deux" {
		log.Fatal(m)
	}
}
func Test_structToMapExclude(t *testing.T) {
	tx := &Tx{}
	tx.Mapper = reflectx.NewMapperFunc("db", CamelToSnake)
	type S struct {
		Ok     int
		Coral  string
		UnDeux string
	}
	m := tx.structToMapExclude(&S{2, "coral", "un_deux"}, []string{"coral"})
	if m["ok"].(int) != 2 {
		log.Fatal(m)
	}
	if _, ok := m["coral"]; ok == true {
		log.Fatal(m)
	}
	if m["un_deux"].(string) != "un_deux" {
		log.Fatal(m)
	}
}
