package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	connString string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOSTNAME"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	connectionPool *sql.DB = nil
)

func init() {
	createConnectionPool()
}

func createConnectionPool() {
	c, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalln("database.go: createConnectionPool: error creating connection pool")
	}
	// unlimited open connection
	c.SetMaxOpenConns(0)
	// max idle connections
	c.SetMaxIdleConns(42)
	connectionPool = c
}
