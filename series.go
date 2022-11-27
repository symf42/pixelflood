package main

import (
	"database/sql"
	"fmt"
)

type Series struct {
	Id          int            `json:"-"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CreatedAt   sql.NullTime   `json:"createdAt"`
	UpdatedAt   sql.NullTime   `json:"updatedAt"`
	DeletedAt   sql.NullTime   `json:"deletedAt"`
	Seasons     map[int]Season `json:"seasons"`
}

type SeriesBulk struct {
	Series struct {
		Id          int
		Title       string
		Description string
		CreatedAt   sql.NullTime
		UpdatedAt   sql.NullTime
		DeletedAt   sql.NullTime
	}
	Season struct {
		Id          sql.NullInt64
		SeriesId    sql.NullInt64
		Number      sql.NullInt64
		Title       sql.NullString
		Description sql.NullString
		CreatedAt   sql.NullTime
		UpdatedAt   sql.NullTime
		DeletedAt   sql.NullTime
		Episodes    []Episode
	}
	Episode struct {
		Id          sql.NullInt64
		SeasonId    sql.NullInt64
		Number      int
		Title       string
		Description string
		CreatedAt   sql.NullTime
		UpdatedAt   sql.NullTime
		DeletedAt   sql.NullTime
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

func getAllSeriesBulk(conn *sql.DB) ([]SeriesBulk, error) {

	series := make([]SeriesBulk, 0)

	stmt, err := conn.Prepare(`SELECT
s.id, s.title, s.description, s.created_at, s.updated_at, s.deleted_at,
se.id, se.series_id, se.number, se.title, se.description, se.created_at, se.updated_at, se.deleted_at,
e.id, e.season_id, e.number, e.title, e.description, e.created_at, e.updated_at, e.deleted_at,
f.id, f.format_id, f.language_id, f.movie_id, f.episode_id, f.hash, f.created_at, f.updated_at, f.deleted_at,
fo.id, fo.name,
l.id, l.short, l.name
FROM series AS s
LEFT JOIN season AS se ON se.series_id = s.id
LEFT JOIN episode AS e ON e.season_id = se.id
LEFT JOIN file AS f ON f.episode_id = e.id
LEFT JOIN format AS fo ON fo.id = f.format_id
LEFT JOIN language AS l ON l.id = f.language_id
ORDER BY s.id, se.id, e.number;`)
	if err != nil {
		return nil, fmt.Errorf("series.go: getAllSeriesBulk(): error preparing statment")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("series.go: getAllSeriesBulk(): error executing query")
	}

	for rows.Next() {

		s := SeriesBulk{}

		if err := rows.Scan(
			&s.Series.Id, &s.Series.Title, &s.Series.Description, &s.Series.CreatedAt, &s.Series.UpdatedAt, &s.Series.DeletedAt,
			&s.Season.Id, &s.Season.SeriesId, &s.Season.Number, &s.Season.Title, &s.Season.Description, &s.Season.CreatedAt, &s.Season.UpdatedAt, &s.Season.DeletedAt,
			&s.Episode.Id, &s.Episode.SeasonId, &s.Episode.Number, &s.Episode.Title, &s.Episode.Description, &s.Episode.CreatedAt, &s.Episode.UpdatedAt, &s.Episode.DeletedAt,
			&s.File.Id, &s.File.FormatId, &s.File.LanguageId, &s.File.MovieId, &s.File.EpisodeId, &s.File.Hash, &s.File.CreatedAt, &s.File.UpdatedAt, &s.File.DeletedAt,
			&s.Format.Id, &s.Format.Name,
			&s.Language.Id, &s.Language.Short, &s.Language.Name,
		); err != nil {
			return nil, fmt.Errorf("series.go: getAllSeriesBulk(): error scanning row: %s", err.Error())
		}

		series = append(series, s)

	}

	return series, nil

}
