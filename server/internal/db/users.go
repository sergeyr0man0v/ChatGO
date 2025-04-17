package db

import (
	"context"
	"fmt"
	"server/internal/models"
)

func (r *repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := "INSERT INTO users(username, encrypted_password) VALUES ($1, $2) returning id, created_at"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.EncryptedPassword).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	u := models.User{}
	query :=
		"SELECT id, username, encrypted_password, created_at, last_login, status FROM users WHERE username = $1"

	// ПРОВЕРИТЬ КОРРЕКТНОСТЬ!
	err := r.db.QueryRowContext(ctx, query, username).Scan(&u.ID, &u.Username, &u.EncryptedPassword, &u.CreatedAt, &u.LastLogin, &u.Status)
	if err != nil {
		return nil, err
	}

	fmt.Println("after error")

	return &u, nil
}

// GetUserByID получает пользователя по его ID
func (r *repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query :=
		"SELECT id, username, encrypted_password, created_at, last_login, status FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username,
		&user.EncryptedPassword, &user.CreatedAt, &user.LastLogin, &user.Status)
	return &user, err
}

// GetAllUsers возвращает всех пользователей
func (r *repository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, username, encrypted_password, created_at, last_login, status FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username,
			&user.EncryptedPassword,
			&user.CreatedAt, &user.LastLogin, &user.Status); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// // InsertUser вставляет нового пользователя в базу данных
// func (db *Database) InsertUser(ctx context.Context, user models.User) error {
// 	query := `INSERT INTO users (id, username, encrypted_password, created_at, last_login, status)
//         	VALUES ($1, $2, $3, $4, $5, $6)`
// 	_, err := db.GetDB().QueryContext(ctx, query, user.ID, user.Username, user.EncryptedPassword,
// 		user.CreatedAt, user.LastLogin, user.Status)
// 	return err
// }
