package db

import (
	"chatgo/server/internal/models"
	"context"
)

// CreateMessage добавляет новое сообщение в базу данных, устанавливает created_at и updated_at CURRENT_TIMESTAMP
func (r *repository) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	query := `
		INSERT INTO messages (
			sender_id, 
			chat_room_id, 
			encrypted_content,
			created_at,
			updated_at,
			is_edited
		) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false)
		RETURNING id, sender_id, chat_room_id, encrypted_content, created_at, updated_at, is_edited`

	err := r.db.QueryRowContext(
		ctx,
		query,
		message.SenderID,
		message.ChatRoomID,
		message.EncryptedContent,
	).Scan(
		&message.ID,
		&message.SenderID,
		&message.ChatRoomID,
		&message.EncryptedContent,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.IsEdited,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessagesByChatRoomID получает сообщения по ID чата, используя лимит для ограничения числа возвращаемых сообщений
func (r *repository) GetMessagesByChatRoomID(ctx context.Context, chatRoomID string, limit int) ([]*models.Message, error) {
	query := `
		SELECT 
			id,
			sender_id,
			chat_room_id,
			encrypted_content,
			created_at,
			updated_at,
			is_edited 
		FROM messages 
		WHERE chat_room_id = $1
		ORDER BY created_at ASC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, chatRoomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		if err = rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ChatRoomID,
			&message.EncryptedContent,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.IsEdited,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// GetMessageByID получает сообщение по ID сообщения
func (r *repository) GetMessageByID(ctx context.Context, messageID string) (*models.Message, error) {
	query := `
		SELECT 
			id,
			sender_id,
			chat_room_id,
			encrypted_content,
			created_at,
			updated_at,
			is_edited 
		FROM messages 
		WHERE id = $1`

	message := &models.Message{}
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(
		&message.ID,
		&message.SenderID,
		&message.ChatRoomID,
		&message.EncryptedContent,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.IsEdited,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}
