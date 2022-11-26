package main

import (
	"database/sql"
	"fmt"
)

type Series struct {
	Id          int          `json:"-"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatedAt   sql.NullTime `json:"createdAt"`
	UpdatedAt   sql.NullTime `json:"updatedAt"`
	DeletedAt   sql.NullTime `json:"deletedAt"`
	Seasons     []Season     `json:"seasons"`
}

func getAllSeries(conn *sql.DB) ([]Series, error) {

	series := make([]Series, 0)

	stmt, err := conn.Prepare("SELECT s.id, s.title, s.description, s.created_at, s.updated_at, s.deleted_at FROM series AS s;")
	if err != nil {
		return nil, fmt.Errorf("series.go: getAllSeries(): error preparing statment")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("series.go: getAllSeries(): error executing query")
	}

	for rows.Next() {

		s := Series{}

		if err := rows.Scan(&s.Id, &s.Title, &s.Description, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt); err != nil {
			return nil, fmt.Errorf("series.go: getAllSeries(): error scanning row: %s", err.Error())
		}

		s.Seasons, err = getSeasonsForSeriesId(s.Id, conn)
		if err != nil {
			return nil, fmt.Errorf("series.go: getAllSeries(): error receiving seasons for series with id '%d': %s", s.Id, err.Error())
		}

		series = append(series, s)

	}

	return series, nil

}
