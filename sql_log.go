package webo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var DB_PG = 0
var DB_ACCESS = 1
var DB_MSSQL = 2

var sql_log_re_question = regexp.MustCompile(`\?`)
var sql_log_re_dollar = regexp.MustCompile(`\$\d+`)

func sql_fake(db_type int, query string, args ...interface{}) string {
	if len(args) == 0 {
		return query
	}
	rqi := 0

	frq := func(s string) string {
		if s == "?" {
			rqi++
			return sql_quoter(db_type, args[rqi-1])
		}
		rqi, _ = strconv.Atoi(s[1:len(s)])
		return sql_quoter(db_type, args[rqi-1])
	}
	switch db_type {
	case DB_ACCESS:
		return sql_log_re_question.ReplaceAllStringFunc(query, frq)
	case DB_PG:
		return sql_log_re_dollar.ReplaceAllStringFunc(query, frq)
	}
	return query
}

func sql_quoter(db_type int, s interface{}) string {
	switch v := s.(type) {
	case nil:
		return "null"
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return "'" + strings.Replace(v, "'", "''", -1) + "'"
	case time.Time:
		return v.Format("'2006-01-02 15:04:05'")
	case *time.Time:
		return v.Format("'2006-01-02 15:04:05'")
	case bool:
		switch db_type {
		case DB_ACCESS:
			switch v {
			case true:
				return "-1"
			case false:
				return "0"
			}
		default:
			switch v {
			case true:
				return "true"
			case false:
				return "false"
			}
		}
	}
	return fmt.Sprintf("%s", s)
}
