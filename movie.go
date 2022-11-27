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
	Files            map[int]File  `json:"files"`
}

type MovieBulk struct {
	Movie struct {
		Id               int
		Title            string
		Description      string
		PrecedingMovieId sql.NullInt64
		CreatedAt        sql.NullTime
		UpdatedAt        sql.NullTime
		DeletedAt        sql.NullTime
	}
	File struct {
		Id         sql.NullInt64
		FormatId   sql.NullInt64
		LanguageId sql.NullInt64
		MovieId    sql.NullInt64
		EpisodeId  sql.NullInt64
		Hash       sql.NullString
		CreatedAt  sql.NullTime
		UpdatedAt  sql.NullTime
		DeletedAt  sql.NullTime
	}
	Format struct {
		Id   sql.NullInt64
		Name sql.NullString
	}
	Language struct {
		Id    sql.NullInt64
		Short sql.NullString
		Name  sql.NullString
	}
}

func getAllMoviesBulk(conn *sql.DB) ([]MovieBulk, error) {

	movies := make([]MovieBulk, 0)

	stmt, err := conn.Prepare(`SELECT
m.id, m.title, m.description, m.preceding_movie_id, m.created_at, m.updated_at, m.deleted_at,
f.id, f.format_id, f.language_id, f.movie_id, f.episode_id, f.hash, f.created_at, f.updated_at, f.deleted_at,
fo.id, fo.name,
l.id, l.short, l.name
FROM movie AS m
LEFT JOIN file AS f ON f.movie_id = m.id
LEFT JOIN format AS fo ON fo.id = f.format_id
LEFT JOIN language AS l ON l.id = f.language_id
ORDER BY m.id ASC;`)
	if err != nil {
		return nil, fmt.Errorf("movie.go: getAllMoviesBulk(): error preparing statment")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("movie.go: getAllMoviesBulk(): error executing query")
	}

	for rows.Next() {

		m := MovieBulk{}

		if err := rows.Scan(
			&m.Movie.Id, &m.Movie.Title, &m.Movie.Description, &m.Movie.PrecedingMovieId, &m.Movie.CreatedAt, &m.Movie.UpdatedAt, &m.Movie.DeletedAt,
			&m.File.Id, &m.File.FormatId, &m.File.LanguageId, &m.File.MovieId, &m.File.EpisodeId, &m.File.Hash, &m.File.CreatedAt, &m.File.UpdatedAt, &m.File.DeletedAt,
			&m.Format.Id, &m.Format.Name,
			&m.Language.Id, &m.Language.Short, &m.Language.Name,
		); err != nil {
			return nil, fmt.Errorf("movie.go: getAllMoviesBulk(): error scanning row: %s", err.Error())
		}

		movies = append(movies, m)

	}

	return movies, nil

}
