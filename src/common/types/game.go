package types

import (
	"database/sql"
	"time"
)

type Game struct {
	Id         Uuid         `json:"id"`
	FinishedAt sql.NullTime `json:"finished_at"`
	CreatedAt  time.Time    `json:"created_at"`
}

type User struct {
	Id        Uuid      `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
