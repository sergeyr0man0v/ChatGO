package models

import (
	"database/sql"
	"time"
)

// UserStatus представляет собой тип для состояния пользователя
type UserStatus int

const (
	Online UserStatus = iota
	Offline
	Away
	Banned
)

// User представляет собой модель пользователя
type User struct {
	ID                int          `json:"id"`
	Username          string       `json:"username"`
	EncryptedPassword string       `json:"encrypted_password"`
	CreatedAt         time.Time    `json:"created_at"`
	LastLogin         sql.NullTime `json:"last_login"` // может быть NULL
	Status            UserStatus   `json:"status"`
}
