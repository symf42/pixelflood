package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id             int              `json:"id"`
	Firstname      string           `json:"-"`
	Lastname       string           `json:"-"`
	Email          string           `json:"email"`
	HashedPassword string           `json:"-"`
	ActivatedAt    sql.NullTime     `json:"-"`
	CreatedAt      sql.NullTime     `json:"-"`
	UpdatedAt      sql.NullTime     `json:"-"`
	DeletedAt      sql.NullTime     `json:"-"`
	Groups         map[string]Group `json:"groups"`
	basicAuth      string
}

type Group struct {
	Id          int                   `json:"-"`
	Name        string                `json:"groupName"`
	Permissions map[string]Permission `json:"permissions"`
}

type Permission struct {
	Id   int    `json:"-"`
	Name string `json:"permissionName"`
}

type UserBulk struct {
	User struct {
		Id             int
		Firstname      string
		Lastname       string
		Email          string
		HashedPassword string
		ActivatedAt    sql.NullTime
		CreatedAt      sql.NullTime
		UpdatedAt      sql.NullTime
		DeletedAt      sql.NullTime
	}
	Group struct {
		Id   int
		Name string
	}
	Permission struct {
		Id   int
		Name string
	}
}

func getAllUsers(conn *sql.DB) ([]UserBulk, error) {

	users := make([]UserBulk, 0)

	stmt, err := conn.Prepare(`SELECT 
u.id, u.firstname, u.lastname, u.email, u.password, u.activated_at, u.created_at, u.updated_at, u.deleted_at,
g.id, g.name, 
p.id, p.name
FROM user AS u
LEFT JOIN user_group AS ug ON ug.user_id = u.id
LEFT JOIN ` + "`" + `group` + "`" + ` AS g ON g.id = ug.group_id
LEFT JOIN group_permission AS gp ON gp.group_id = g.id
LEFT JOIN permission AS p ON p.id = gp.permission_id;`)
	if err != nil {
		return nil, fmt.Errorf("user.go: getAllUsers(): error preparing statment")
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("user.go: getAllUsers(): error executing query")
	}

	for rows.Next() {

		u := UserBulk{}

		if err := rows.Scan(
			&u.User.Id, &u.User.Firstname, &u.User.Lastname, &u.User.Email, &u.User.HashedPassword, &u.User.ActivatedAt, &u.User.CreatedAt, &u.User.UpdatedAt, &u.User.DeletedAt,
			&u.Group.Id, &u.Group.Name,
			&u.Permission.Id, &u.Permission.Name,
		); err != nil {
			return nil, fmt.Errorf("user.go: getAllUsers(): error scanning row: %s", err.Error())
		}

		users = append(users, u)

	}

	return users, nil
}
