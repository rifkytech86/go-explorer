package dto

import "database/sql"

type Users struct {
	ID        int64        `json:"diary_id"`
	Email     string       `json:"user_id"`
	Username  string       `json:"username"`
	Password  string       `json:"password"`
	Name      string       `json:"diary_desc"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
