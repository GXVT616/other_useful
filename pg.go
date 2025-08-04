package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//===

type pg_conn struct {
	host     string
	port     int64
	user     string
	password string
	dbname   string
	connect  *sql.DB
	limit    int64
	offset   int64
}

//---

func (p *pg_conn) Connect() *pg_conn {

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		(*p).host, (*p).port, (*p).user, (*p).password, (*p).dbname)
	var err error
	(*p).connect, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	err = (*p).connect.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	return p
}

//---

func (p *pg_conn) Query(query string) ([][]any, bool) {

	result := make([][]any, 0)

	if (*p).limit != 0 {
		query = fmt.Sprintf("%v LIMIT %v OFFSET %v;", query, (*p).limit, (*p).offset)
		// if (*p).count == 0{

		// }
	}
	(*p).offset += (*p).limit
	//query := fmt.Sprintf("SELECT mcast_main,source_m_mcast_main FROM (SELECT  mcast_main,source_m_mcast_main FROM info_mcast_dth_center UNION SELECT  mcast_main,source_m_mcast_main FROM info_mcast_dth_siberia) as f ORDER BY mcast_main,source_m_mcast_main LIMIT %v OFFSET %v;", (*p).limit, (*p).offset)

	rows, err := (*p).connect.Query(query)
	if err != nil {
		log.Println("[Postgres] Ошибка с получением данных:", err)
		return result, false
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Fatalln("[Postgres] Ошибка с закрытием строки:", err)
		}
	}(rows)

	count_column := func() int64 {
		column, err := rows.Columns()
		if err != nil {
			log.Println("[PG] Проблема с определением колонок в ответе", err)
		}
		return int64(len(column))
	}()

	for rows.Next() {

		rows_value := make([]any, 0)
		for range count_column {
			var val any
			rows_value = append(rows_value, &val)
		}
		rows.Scan(rows_value...)

		for i, v := range rows_value {

			switch v.(type) {
			case *any:
				rows_value[i] = *(v.(*any))
			}
		}

		// for i, v := range rows_value {

		// 	switch rows_value[i].(type) {
		// 	case []uint8:
		// 		rows_value[i] = string(v.([]uint8))

		// 	}
		// }

		result = append(result, rows_value)
	}

	if len(result) == 0 {
		return result, false
	}

	return result, true
}
func (p *pg_conn) Write(query string, arg ...any) {

	rows, err := (*p).connect.Query(query, arg...)
	if err != nil {
		log.Println("[Postgres] Ошибка с получением данных:", err)
		return
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Fatalln("[Postgres] Ошибка с закрытием строки:", err)
		}
	}(rows)

}

func Scan(data []any, input_type ...any) {
	for i, v := range input_type {
		switch v.(type) {
		case *string:
			if new_val, ok := data[i].(string); ok {
				*(input_type[i].(*string)) = new_val
			}
		case *int:
			if new_val, ok := data[i].(int); ok {
				*(input_type[i].(*int)) = new_val
			}
		case *int64:
			if new_val, ok := data[i].(int64); ok {
				*(input_type[i].(*int64)) = new_val
			}
		case *[]uint8:
			if new_val, ok := data[i].([]uint8); ok {
				*(input_type[i].(*[]uint8)) = new_val
			}
		default:
			log.Printf("не могу преобразовать - %T", v)
		}

	}
}

//===
