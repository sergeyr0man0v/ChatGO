package models

import (
	"database/sql"
	"time"
)

// UserStatus представляет собой тип для состояния пользователя
type UserStatus string

const (
	Online  string = "online"
	Offline string = "offline"
	Away    string = "away"
	Banned  string = "banned"
)

// User представляет собой модель пользователя
type User struct {
	ID                string       `json:"id"`
	Username          string       `json:"username"`
	EncryptedPassword string       `json:"encrypted_password"`
	CreatedAt         time.Time    `json:"created_at"`
	LastLogin         sql.NullTime `json:"last_login"` // может быть NULL
	Status            UserStatus   `json:"status"`
}
