package db

import (
	"chat/internal/models"
)

// InsertUser вставляет нового пользователя в базу данных
func (db *Database) InsertUser(user models.User) error {
	query := `INSERT INTO users (id, username, encrypted_password, created_at, last_login, status) 
        	VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.GetDB().Exec(query, user.ID, user.Username, user.EncryptedPassword,
		user.CreatedAt, user.LastLogin, user.Status)
	return err
}

// GetUserByID получает пользователя по его ID
func (db *Database) GetUserByID(id string) (models.User, error) {
	var user models.User
	query :=
		"SELECT id, username, encrypted_password, created_at, last_login, status FROM users WHERE id = $1"
	err := db.GetDB().QueryRow(query, id).Scan(&user.ID, &user.Username,
		&user.EncryptedPassword, &user.CreatedAt, &user.LastLogin, &user.Status)
	return user, err
}

// GetAllUsers возвращает всех пользователей
func (db *Database) GetAllUsers() ([]models.User, error) {
	rows, err := db.GetDB().Query(
		"SELECT id, username, encrypted_password, created_at, last_login, status FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username,
			&user.EncryptedPassword,
			&user.CreatedAt, &user.LastLogin, &user.Status); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
