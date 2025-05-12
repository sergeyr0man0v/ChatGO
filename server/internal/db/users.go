package db

import (
	"context"
	"chatgo/server/internal/models"
)

// CreateUser добавляет нового пользователя в базу данных, устанавливает created_at и last_login CURRENT_TIMESTAMP
func (r *repository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users(
			username, 
			encrypted_password,
			created_at,
			last_login,
			status
		) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $3) 
		RETURNING id, username, encrypted_password, created_at, last_login, status`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.EncryptedPassword,
		user.Status,
	).Scan(
		&user.ID,
		&user.Username,
		&user.EncryptedPassword,
		&user.CreatedAt,
		&user.LastLogin,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername получает пользователя по его username
func (r *repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.User{}
	query := `
		SELECT 
			id, 
			username, 
			encrypted_password, 
			created_at, 
			last_login, 
			status 
		FROM users 
		WHERE username = $1`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.EncryptedPassword,
		&user.CreatedAt,
		&user.LastLogin,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID получает пользователя по его ID
func (r *repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `
		SELECT 
			id, 
			username, 
			encrypted_password, 
			created_at, 
			last_login, 
			status 
		FROM users 
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.EncryptedPassword,
		&user.CreatedAt,
		&user.LastLogin,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// GetAllUsers возвращает всех пользователей
func (r *repository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT 
			id, 
			username, 
			encrypted_password, 
			created_at, 
			last_login, 
			status 
		FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.EncryptedPassword,
			&user.CreatedAt,
			&user.LastLogin,
			&user.Status,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
