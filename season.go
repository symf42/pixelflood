package main

import (
	"database/sql"
	"fmt"
)

type Season struct {
	Id          int           `json:"-"`
	SeriesId    sql.NullInt64 `json:"seriesId"`
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   sql.NullTime  `json:"createdAt"`
	UpdatedAt   sql.NullTime  `json:"updatedAt"`
	DeletedAt   sql.NullTime  `json:"deletedAt"`
	Episodes    []Episode     `json:"episodes"`
}

func getSeasonsForSeriesId(id int, conn *sql.DB) ([]Season, error) {

	seasons := make([]Season, 0)

	stmt, err := conn.Prepare("SELECT s.id, s.series_id, s.number, s.title, s.description, s.created_at, s.updated_at, s.deleted_At FROM season AS s WHERE s.series_id = ?;")
	if err != nil {
		return nil, fmt.Errorf("season.go: getSeasonsForSeriesId(): error preparing statment")
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("season.go: getSeasonsForSeriesId(): error executing query")
	}

	for rows.Next() {

		s := Season{}

		if err := rows.Scan(&s.Id, &s.SeriesId, &s.Number, &s.Title, &s.Description, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt); err != nil {
			return nil, fmt.Errorf("season.go: getSeasonsForSeriesId(): error scanning row: %s", err.Error())
		}

		s.Episodes, err = getEpisodesForSeasonId(s.Id, conn)
		if err != nil {
			return nil, fmt.Errorf("season.go: getSeasonsForSeriesId(): error receiving episodes for season with id '%d': %s", s.Id, err.Error())
		}

		seasons = append(seasons, s)

	}

	return seasons, nil

}
