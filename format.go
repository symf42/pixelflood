package main

import (
	"database/sql"
	"fmt"
)

type Format struct {
	Id   int    `json:"-"`
	Name string `json:"name"`
}

func getFormatById(id int, conn *sql.DB) (*Format, error) {

	f := Format{}

	stmt, err := conn.Prepare("SELECT f.id, f.name FROM format AS f WHERE id = ?;")
	if err != nil {
		return nil, fmt.Errorf("format.go: getFormatById(): error preparing statment")
	}

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		return nil, fmt.Errorf("format.go: getFormatById(): error executing query")
	}

	if err := row.Scan(&f.Id, &f.Name); err != nil {
		return nil, fmt.Errorf("format.go: getFormatById(): error scanning row: %s", err.Error())
	}

	return &f, nil

}
