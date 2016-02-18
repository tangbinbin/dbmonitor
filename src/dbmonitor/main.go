package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	host         = flag.String("h", "127.0.0.1:3306", "mysql addr")
	user         = flag.String("u", "monitor", "user to connect db")
	passwd       = flag.String("p", "monitor", "passwd to connect db")
	connUser     = flag.String("U", "monitor", "user to connect monitor db")
	connPasswd   = flag.String("P", "monitor", "passwd to connect monitor db")
	timeInterval = flag.Uint64("t", 10, "time interval to monitor s")
	db           *sql.DB
	servers      DbServer
)

type Info struct {
	questions    uint64
	comSelect    uint64
	comInsert    uint64
	comUpdate    uint64
	comDelete    uint64
	rowsInsert   uint64
	rowsSelect   uint64
	rowsUpdate   uint64
	rowsDelete   uint64
	threadCreate uint64
	byteReceived uint64
	byteSent     uint64
}

type Status struct {
	qps  uint64
	sel  uint64
	upd  uint64
	del  uint64
	ins  uint64
	rre  uint64
	rin  uint64
	rup  uint64
	rdel uint64
	con  uint64
	cre  uint64
	run  uint64
	bre  uint64
	bse  uint64
}

type Server struct {
	id      int
	addr    string
	cluster string
	conn    *sql.DB
	info    *Info
	status  *Status
	first   bool
}

type DbServer map[int]*Server

func init() {
	flag.Parse()
	initDb()
	initServers()
}

func initDb() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&timeout=%s",
		*user, *passwd, *host, "dbmonitor", "utf8", "100ms")
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(1)
}

func initServers() {
	ssql := "select id, addr, cluster from db_account"
	rows, err := db.Query(ssql)
	if err != nil {
		log.Fatalf("get db info fail %s %v", ssql, err)
		os.Exit(1)
	}
	defer rows.Close()
	servers = make(map[int]*Server)
	for rows.Next() {
		var (
			id      int
			addr    string
			cluster string
		)
		if err = rows.Scan(&id, &addr, &cluster); err != nil {
			log.Fatal(err)
			return
		}

		connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=100ms",
			*connUser, *connPasswd, addr, "information_schema")
		conn, err := sql.Open("mysql", connStr)
		if err != nil {
			log.Fatal(err)
		}
		servers[id] = &Server{
			id:      id,
			addr:    addr,
			cluster: cluster,
			conn:    conn,
			info:    &Info{},
			status:  &Status{},
			first:   true,
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go monitor()
	<-signals
	db.Close()
}

func monitor() {
	for range time.NewTicker(time.Duration(*timeInterval) * time.Second).C {
		for _, server := range servers {
			go getInfo(server)
		}
	}
}

func getInfo(s *Server) {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}
	}()
	ssql := "show global status where variable_name in " +
		"('Questions', 'Com_select', 'Com_update', " +
		"'Com_insert', 'Com_delete', 'Threads_connected', " +
		"'Threads_created', 'Threads_running', " +
		"'Innodb_rows_inserted', 'Innodb_rows_read', " +
		"'Innodb_rows_updated', 'Innodb_rows_deleted', " +
		"'Bytes_received','Bytes_sent')"
	rows, e := s.conn.Query(ssql)
	if e != nil {
		panic(e)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			key   string
			value uint64
		)
		if err := rows.Scan(&key, &value); err != nil {
			return
		}
		switch key {
		case "Bytes_received":
			s.status.bre = (value - s.info.byteReceived) / *timeInterval
			s.info.byteReceived = value
		case "Bytes_sent":
			s.status.bse = (value - s.info.byteSent) / *timeInterval
			s.info.byteSent = value
		case "Com_delete":
			s.status.del = (value - s.info.comDelete) / *timeInterval
			s.info.comDelete = value
		case "Questions":
			s.status.qps = (value - s.info.questions) / *timeInterval
			s.info.questions = value
		case "Com_select":
			s.status.sel = (value - s.info.comSelect) / *timeInterval
			s.info.comSelect = value
		case "Com_update":
			s.status.upd = (value - s.info.comUpdate) / *timeInterval
			s.info.comUpdate = value
		case "Com_insert":
			s.status.ins = (value - s.info.comInsert) / *timeInterval
			s.info.comInsert = value
		case "Innodb_rows_inserted":
			s.status.rin = (value - s.info.rowsInsert) / *timeInterval
			s.info.rowsInsert = value
		case "Innodb_rows_read":
			s.status.rre = (value - s.info.rowsSelect) / *timeInterval
			s.info.rowsSelect = value
		case "Innodb_rows_updated":
			s.status.rup = (value - s.info.rowsUpdate) / *timeInterval
			s.info.rowsUpdate = value
		case "Innodb_rows_deleted":
			s.status.rdel = (value - s.info.rowsDelete) / *timeInterval
			s.info.rowsDelete = value
		case "Threads_connected":
			s.status.con = value
		case "Threads_created":
			s.status.cre = (value - s.info.threadCreate) / *timeInterval
			s.info.threadCreate = value
		case "Threads_running":
			s.status.run = value
		default:
		}
	}
	if s.first {
		s.first = false
		return
	}
	saveStatus(s)
	checkStatus(s)
}

func saveStatus(s *Server) {
	field := []string{
		"id",
		"qps",
		"com_select",
		"com_delete",
		"com_update",
		"com_insert",
		"rows_read",
		"rows_delete",
		"rows_insert",
		"rows_update",
		"thread_created",
		"thread_connected",
		"thread_running",
		"byte_received",
		"byte_sent",
		"created_time",
	}
	values := []string{}
	for i := 1; i <= len(field); i++ {
		values = append(values, "?")
	}
	ssql := "insert into db_status (" +
		strings.Join(field, ",") +
		") values (" +
		strings.Join(values, ",") + ")"
	_, err := db.Exec(ssql,
		s.id,
		s.status.qps,
		s.status.sel,
		s.status.del,
		s.status.upd,
		s.status.ins,
		s.status.rre,
		s.status.rdel,
		s.status.rin,
		s.status.rup,
		s.status.cre,
		s.status.con,
		s.status.run,
		s.status.bre,
		s.status.bse,
		time.Now().Unix())
	if err != nil {
		log.Println(err)
	}
}

func checkStatus(s *Server) {
	field := []string{
		"qps",
		"com_select",
		"com_delete",
		"com_update",
		"com_insert",
		"rows_read",
		"rows_delete",
		"rows_insert",
		"rows_update",
		"thread_created",
		"thread_connected",
		"thread_running",
		"byte_received",
		"byte_sent",
	}
	ssql := "select " + strings.Join(field, ",") +
		" from db_shreshold where id = ?"
	rows, _ := db.Query(ssql, s.id)
	var (
		qps, sel, del, upd, ins uint64
		rre, rdel, rin, rup     uint64
		cre, con, run, bre, bse uint64
	)
	for rows.Next() {
		if err := rows.Scan(&qps, &sel, &del, &upd, &ins,
			&rre, &rdel, &rin, &rup, &cre, &con, &run,
			&bre, &bse); err != nil {
			log.Println(err)
			return
		}
	}
	if s.status.qps > qps {
		log.Println("qps over")
	}
	if s.status.sel > sel {
		log.Println("com_select over")
	}
	if s.status.del > del {
		log.Println("com_delete over")
	}
	if s.status.upd > upd {
		log.Println("com_update over")
	}
	if s.status.ins > ins {
		log.Println("com_insert over")
	}
	if s.status.rre > rre {
		log.Println("rows_read over")
	}
	if s.status.rdel > rdel {
		log.Println("rows_delete over")
	}
	if s.status.rin > rin {
		log.Println("rows_insert over")
	}
	if s.status.rup > rup {
		log.Println("rows_update over")
	}
	if s.status.cre > cre {
		log.Println("tread create over")
	}
	if s.status.con > con {
		log.Println("conn over")
	}
	if s.status.run > run {
		log.Println("thread running over")
	}
	if s.status.bre > bre {
		log.Println("byte_received over")
	}
	if s.status.bse > bse {
		log.Println("byte_sent over")
	}
}
