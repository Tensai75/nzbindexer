package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var mysqlError *mysql.MySQLError

func connectMySQL() error {

	var err error
	db, err = sql.Open("mysql", conf.Database.User+":"+conf.Database.Password+"@tcp("+conf.Database.Host+":"+strconv.Itoa(conf.Database.Port)+")/"+conf.Database.Name)
	if err != nil {
		fmt.Printf("Connection to database failed: %v\n", err)
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(conf.ParallelScans)
	db.SetMaxIdleConns(conf.ParallelScans)

	if err := db.Ping(); err != nil {
		fmt.Printf("Unable to reach database: %v\n", err)
		return err
	}
	fmt.Println("Connection to database established")

	return nil

}
