package db

import "chat/internal/models"

// InsertMessage вставляет новое сообщение в базу данных
func (db *Database) InsertMessage(message models.Message) error {
	query := `INSERT INTO messages (id, sender_id, chat_room_id, encrypted_content,
			reply_to_message_id, created_at, updated_at, is_edited) 
        	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.GetDB().Exec(query, message.ID, message.SenderID, message.ChatRoomID,
		message.EncryptedContent, message.ReplyToMessageID, message.CreatedAt,
		message.UpdatedAt, message.IsEdited)
	return err
}

// GetMessagesByChatRoom получает сообщения по ID чата
func (db *Database) GetMessagesByChatRoom(chatRoomID string) ([]models.Message, error) {
	rows, err := db.GetDB().Query(
		`SELECT id, sender_id, chat_room_id, encrypted_content, reply_to_message_id,
		created_at, updated_at, is_edited FROM messages WHERE chat_room_id = $1`, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.SenderID,
			&message.ChatRoomID, &message.EncryptedContent,
			&message.ReplyToMessageID, &message.CreatedAt,
			&message.UpdatedAt, &message.IsEdited); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// GetMessageByID получает сообщение по его ID
func (db *Database) GetMessageByID(id string) (models.Message, error) {
	var message models.Message
	query := `SELECT id, sender_id, chat_room_id, encrypted_content, reply_to_message_id,
				created_at, updated_at, is_edited FROM messages WHERE id = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&message.ID, &message.SenderID,
		&message.ChatRoomID, &message.EncryptedContent, &message.ReplyToMessageID,
		&message.CreatedAt, &message.UpdatedAt, &message.IsEdited)
	return message, err
}
