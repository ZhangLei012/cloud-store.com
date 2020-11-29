package conn

import (
	"database/sql"
	"log"

	"cloud-store.com/service/dbproxy/config"
)

var db *sql.DB

//InitDBConn 初始化数据库连接
func InitDBConn() {
	var err error
	db, err = sql.Open("mysql", config.MySQLSource)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	cloumns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(cloumns))
	values := make([]interface{}, len(cloumns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		for i, value := range values {
			if value != nil {
				record[cloumns[i]] = value
			}
		}
		records = append(records, record)
	}
	return records
}
