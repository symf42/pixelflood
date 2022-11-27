package main

import (
	"database/sql"
	"fmt"
)

type File struct {
	Id         int            `json:"-"`
	FormatId   sql.NullInt64  `json:"-"`
	LanguageId sql.NullInt64  `json:"-"`
	MovieId    sql.NullInt64  `json:"-"`
	EpisodeId  sql.NullInt64  `json:"-"`
	Hash       sql.NullString `json:"hash"`
	CreatedAt  sql.NullTime   `json:"-"`
	UpdatedAt  sql.NullTime   `json:"-"`
	DeletedAt  sql.NullTime   `json:"-"`
	Format     Format         `json:"format"`
	Language   Language       `json:"language"`
}

func getFilesForMovieId(movieId int, conn *sql.DB) ([]File, error) {

	files := make([]File, 0)

	stmt, err := conn.Prepare("SELECT f.id, f.format_id, f.language_id, f.movie_id, f.episode_id, f.hash, f.created_at, f.updated_at, f.deleted_at FROM file AS f WHERE f.movie_id = ?;")
	if err != nil {
		return nil, fmt.Errorf("file.go: getFilesForMovieId(): error preparing statment")
	}

	rows, err := stmt.Query(movieId)
	if err != nil {
		return nil, fmt.Errorf("file.go: getFilesForMovieId(): error executing query")
	}

	for rows.Next() {

		f := File{}

		if err := rows.Scan(&f.Id, &f.FormatId, &f.LanguageId, &f.MovieId, &f.EpisodeId, &f.Hash, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
			return nil, fmt.Errorf("file.go: getFilesForMovieId(): error scanning row: %s", err.Error())
		}

		if f.FormatId.Valid {
			fo, err := getFormatById(int(f.FormatId.Int64), conn)
			if err != nil {
				return nil, fmt.Errorf("file.go: getFilesForMovieId(): error loading format for formatId: %d with error: %s", f.FormatId.Int64, err.Error())
			}
			f.Format = *fo
		}

		if f.LanguageId.Valid {
			l, err := getLanguageById(int(f.LanguageId.Int64), conn)
			if err != nil {
				return nil, fmt.Errorf("file.go: getFilesForMovieId(): error loading language for languageId: %d with error: %s", f.LanguageId.Int64, err.Error())
			}
			f.Language = *l
		}

		files = append(files, f)

	}

	return files, nil

}

func getFilesForEpisodeId(episodeId int, conn *sql.DB) ([]File, error) {

	files := make([]File, 0)

	stmt, err := conn.Prepare("SELECT f.id, f.format_id, f.language_id, f.movie_id, f.episode_id, f.hash, f.created_at, f.updated_at, f.deleted_at FROM file AS f WHERE f.episode_id = ?;")
	if err != nil {
		return nil, fmt.Errorf("file.go: getFilesForEpisodeId(): error preparing statment")
	}

	rows, err := stmt.Query(episodeId)
	if err != nil {
		return nil, fmt.Errorf("file.go: getFilesForEpisodeId(): error executing query")
	}

	for rows.Next() {

		f := File{}

		if err := rows.Scan(&f.Id, &f.FormatId, &f.LanguageId, &f.MovieId, &f.EpisodeId, &f.Hash, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt); err != nil {
			return nil, fmt.Errorf("file.go: getFilesForEpisodeId(): error scanning row: %s", err.Error())
		}

		if f.FormatId.Valid {
			fo, err := getFormatById(int(f.FormatId.Int64), conn)
			if err != nil {
				return nil, fmt.Errorf("file.go: getFilesForEpisodeId(): error loading format for formatId: %d with error: %s", f.FormatId.Int64, err.Error())
			}
			f.Format = *fo
		}

		if f.LanguageId.Valid {
			l, err := getLanguageById(int(f.LanguageId.Int64), conn)
			if err != nil {
				return nil, fmt.Errorf("file.go: getFilesForEpisodeId(): error loading language for languageId: %d with error: %s", f.LanguageId.Int64, err.Error())
			}
			f.Language = *l
		}

		files = append(files, f)

	}

	return files, nil

}
