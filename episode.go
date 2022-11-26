package main

import (
	"database/sql"
	"fmt"
)

type Episode struct {
	Id          int           `json:"-"`
	SeasonId    sql.NullInt64 `json:"seasonId"`
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   sql.NullTime  `json:"createdAt"`
	UpdatedAt   sql.NullTime  `json:"updatedAt"`
	DeletedAt   sql.NullTime  `json:"deletedAt"`
	Files       []File        `json:"files"`
}

func getEpisodesForSeasonId(id int, conn *sql.DB) ([]Episode, error) {

	episodes := make([]Episode, 0)

	stmt, err := conn.Prepare("SELECT e.id, e.season_id, e.number, e.title, e.description, e.created_at, e.updated_at, e.deleted_at FROM episode AS e WHERE e.season_id = ?;")
	if err != nil {
		return nil, fmt.Errorf("episode.go: getEpisodesForSeasonId(): error preparing statment")
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("episode.go: getEpisodesForSeasonId(): error executing query")
	}

	for rows.Next() {

		e := Episode{}

		if err := rows.Scan(&e.Id, &e.SeasonId, &e.Number, &e.Title, &e.Description, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, fmt.Errorf("episode.go: getEpisodesForSeasonId(): error scanning row: %s", err.Error())
		}

		e.Files, err = getFilesForEpisodeId(e.Id, conn)
		if err != nil {
			return nil, fmt.Errorf("episode.go: getEpisodesForEpisodeId(): error receiving files for episode with id '%d': %s", e.Id, err.Error())
		}

		episodes = append(episodes, e)

	}

	return episodes, nil

}
