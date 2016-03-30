package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

var (
	addr  = flag.String("h", "127.0.0.1:3306", "mysql host")
	ltime = flag.Int("t", 1, "n * hour time")
	db    *sql.DB
)

func init() {
	flag.Parse()
	initDb()
}

func initDb() {
	connStr := "monitor:monitor@tcp(127.0.0.1:3306)/dbmonitor"
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
}

func main() {
	getId := "select id from db_account where addr = ?"
	var id int
	err := db.QueryRow(getId, *addr).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	fields := []string{
		"from_unixtime(created_time)",
		"qps",
		"com_select",
		"com_insert",
		"com_delete",
		"com_update",
	}
	ssql := "select " +
		strings.Join(fields, ",") +
		" from db_status " +
		"where id = ? and " +
		"created_time > (unix_timestamp(now()) - ? * 3600)"
	rows, err := db.Query(ssql, id, *ltime)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var (
			qps  int
			sel  int
			del  int
			upd  int
			ins  int
			time string
		)
		if err = rows.Scan(&time,
			&qps,
			&sel,
			&ins,
			&del,
			&upd,
		); err != nil {
			log.Fatal(err)
		}
		if i%20 == 0 {
			fmt.Println()
			echoHead()
		}
		i++
		fmt.Printf("%s%6d%6d%6d%6d%6d\n",
			time,
			qps,
			sel,
			ins,
			del,
			upd,
		)
	}
}

func echoHead() {
	fmt.Printf("%19s%6s%6s%6s%6s%6s\n",
		"time",
		"qps",
		"sel",
		"ins",
		"del",
		"upd",
	)
}
