package main

import (
	"database/sql"
	"fmt"
)

type Language struct {
	Id    int    `json:"-"`
	Short string `json:"short"`
	Name  string `json:"name"`
}

func getLanguageById(id int, conn *sql.DB) (*Language, error) {

	l := Language{}

	stmt, err := conn.Prepare("SELECT l.id, l.short, l.name FROM language AS l WHERE id = ?;")
	if err != nil {
		return nil, fmt.Errorf("language.go: getLanguageById(): error preparing statment")
	}

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		return nil, fmt.Errorf("language.go: getLanguageById(): error executing query")
	}

	if err := row.Scan(&l.Id, &l.Short, &l.Name); err != nil {
		return nil, fmt.Errorf("language.go: getLanguageById(): error scanning row: %s", err.Error())
	}

	return &l, nil

}
