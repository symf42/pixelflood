package main

import (
	"database/sql"
	"fmt"
)

type Movie struct {
	Id               int           `json:"-"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	PrecedingMovieId sql.NullInt64 `json:"-"`
	CreatedAt        sql.NullTime  `json:"createdAt"`
	UpdatedAt        sql.NullTime  `json:"updatedAt"`
	DeletedAt        sql.NullTime  `json:"deletedAt"`
	Files            []File        `json:"files"`
}

func getAllMovies(conn *sql.DB) ([]Movie, error) {

	movies := make([]Movie, 0)

	stmt, err := conn.Prepare("SELECT m.id, m.title, m.description, m.preceding_movie_id, m.created_at, m.updated_at, m.deleted_at FROM `movie` AS m;")
	if err != nil {
		return nil, fmt.Errorf("movie.go: getAllMovies(): error preparing statment")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("movie.go: getAllMovies(): error executing query")
	}

	for rows.Next() {

		m := Movie{}

		if err := rows.Scan(&m.Id, &m.Title, &m.Description, &m.PrecedingMovieId, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt); err != nil {
			return nil, fmt.Errorf("movie.go: getAllMovies(): error scanning row: %s", err.Error())
		}

		m.Files, err = getFilesForMovieId(m.Id, conn)
		if err != nil {
			return nil, fmt.Errorf("movie.go: getAllMovies(): error receiving files for movie with id '%d': %s", m.Id, err.Error())
		}

		movies = append(movies, m)

	}

	return movies, nil
}
