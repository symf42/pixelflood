package main

import (
	"database/sql"
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
	Files       map[int]File  `json:"files"`
}
