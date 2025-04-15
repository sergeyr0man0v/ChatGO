package db

import "chat/internal/models"

// InsertChatRoom вставляет новый чат в базу данных
func (db *Database) InsertChatRoom(chatRoom models.ChatRoom) error {
	query := `INSERT INTO chat_rooms (id, name, type, created_at, creator_id) 
        	VALUES ($1, $2, $3, $4, $5)`
	_, err := db.GetDB().Exec(query, chatRoom.ID, chatRoom.Name, chatRoom.Type,
		chatRoom.CreatedAt, chatRoom.CreatorID)
	return err
}

// GetChatRoomByID получает чат по его ID
func (db *Database) GetChatRoomByID(id string) (models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	query := "SELECT id, name, type, created_at, creator_id FROM chat_rooms WHERE id = $1"
	err := db.GetDB().QueryRow(query, id).Scan(&chatRoom.ID, &chatRoom.Name,
		&chatRoom.Type, &chatRoom.CreatedAt, &chatRoom.CreatorID)
	return chatRoom, err
}

// GetAllChatRooms возвращает все чаты
func (db *Database) GetAllChatRooms() ([]models.ChatRoom, error) {
	rows, err := db.GetDB().Query("SELECT id, name, type, created_at, creator_id FROM chat_rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatRooms []models.ChatRoom
	for rows.Next() {
		var chatRoom models.ChatRoom
		if err := rows.Scan(&chatRoom.ID, &chatRoom.Name, &chatRoom.Type,
			&chatRoom.CreatedAt, &chatRoom.CreatorID); err != nil {
			return nil, err
		}
		chatRooms = append(chatRooms, chatRoom)
	}
	return chatRooms, nil
}
