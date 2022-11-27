package main

import (
	"database/sql"
)

type Season struct {
	Id          int             `json:"-"`
	SeriesId    sql.NullInt64   `json:"seriesId"`
	Number      int             `json:"number"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CreatedAt   sql.NullTime    `json:"createdAt"`
	UpdatedAt   sql.NullTime    `json:"updatedAt"`
	DeletedAt   sql.NullTime    `json:"deletedAt"`
	Episodes    map[int]Episode `json:"episodes"`
}
