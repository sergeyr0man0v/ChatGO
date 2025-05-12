package db

import (
	"chatgo/server/internal/models"
	"context"
	"database/sql"
)

// CreateChatRoom создает новый чат, усталивая created_at CURRENT_TIMESTAMP.
// А также добавляет создателя в таблицу chat_room_members как админа.
func (r *repository) CreateChatRoom(ctx context.Context, chatRoom *models.ChatRoom) (*models.ChatRoom, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO chat_rooms (name, type, creator_id, created_at) 
        	VALUES ($1, $2, $3, CURRENT_TIMESTAMP) 
			RETURNING id, name, type, creator_id, created_at`

	err = tx.QueryRowContext(ctx, query,
		chatRoom.Name,
		chatRoom.Type,
		chatRoom.CreatorID,
	).Scan(
		&chatRoom.ID,
		&chatRoom.Name,
		&chatRoom.Type,
		&chatRoom.CreatorID,
		&chatRoom.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Добавление админа в таблицу chat_room_members
	member := &models.ChatRoomMember{
		UserID:     chatRoom.CreatorID,
		ChatRoomID: chatRoom.ID,
		MemberRole: models.Admin,
	}

	memberQuery := `INSERT INTO chat_room_members (user_id, chat_room_id, role, joined_at)
					VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
					RETURNING user_id, chat_room_id, role, joined_at`

	err = tx.QueryRowContext(ctx, memberQuery,
		member.UserID,
		member.ChatRoomID,
		member.MemberRole,
	).Scan(
		&member.UserID,
		&member.ChatRoomID,
		&member.MemberRole,
		&member.JoinedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return chatRoom, nil
}

// GetChatRoomByID возвращает чат по ID чата
func (r *repository) GetChatRoomByID(ctx context.Context, chatRoomID string) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	query := `SELECT id, name, type, creator_id, created_at 
			FROM chat_rooms WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, chatRoomID).Scan(
		&chatRoom.ID,
		&chatRoom.Name,
		&chatRoom.Type,
		&chatRoom.CreatorID,
		&chatRoom.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &chatRoom, nil
}

// GetChatRoomsByUserID возвращает все чаты по ID участника
func (r *repository) GetChatRoomsByUserID(ctx context.Context, userID string) ([]*models.ChatRoom, error) {
	query := `SELECT cr.id, cr.name, cr.type, cr.creator_id, cr.created_at,
			FROM chat_rooms cr
			JOIN chat_room_members crm ON cr.id = crm.chat_room_id
			WHERE crm.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatRooms []*models.ChatRoom
	for rows.Next() {
		var chatRoom models.ChatRoom

		err := rows.Scan(
			&chatRoom.ID,
			&chatRoom.Name,
			&chatRoom.Type,
			&chatRoom.CreatorID,
			&chatRoom.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		chatRooms = append(chatRooms, &chatRoom)
	}

	return chatRooms, nil
}

// GetAllChatRooms возвращает все чаты
func (r *repository) GetAllChatRooms(ctx context.Context) ([]*models.ChatRoom, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, type, created_at, creator_id FROM chat_rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatRooms []*models.ChatRoom
	for rows.Next() {
		var chatRoom models.ChatRoom
		if err := rows.Scan(&chatRoom.ID, &chatRoom.Name, &chatRoom.Type,
			&chatRoom.CreatedAt, &chatRoom.CreatorID); err != nil {
			return nil, err
		}
		chatRooms = append(chatRooms, &chatRoom)
	}
	return chatRooms, nil
}

// UpdateChatRoom обновляет имя и тип чата по ID чата
func (r *repository) UpdateChatRoom(ctx context.Context, chatRoom *models.ChatRoom) (*models.ChatRoom, error) {
	query := `UPDATE chat_rooms 
			SET name = $1, type = $2 
			WHERE id = $3 
			RETURNING id, name, type, creator_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		chatRoom.Name,
		chatRoom.Type,
		chatRoom.ID,
	).Scan(
		&chatRoom.ID,
		&chatRoom.Name,
		&chatRoom.Type,
		&chatRoom.CreatorID,
		&chatRoom.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return chatRoom, nil
}

// DeleteChatRoom удаляет чат по ID чата, а также удаляет все записи
// в таблице chat_room_members для этого чата, используя метод [RemoveMembersByChatRoomID]
func (r *repository) DeleteChatRoom(ctx context.Context, chatRoom *models.ChatRoom) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Удаление всех записей в таблице chat_room_members для этого чата
	_, err = tx.ExecContext(ctx, "DELETE FROM chat_room_members WHERE chat_room_id = $1", chatRoom.ID)
	if err != nil {
		return err
	}

	// Удаление чата из таблицы chat_rooms
	_, err = tx.ExecContext(ctx, "DELETE FROM chat_rooms WHERE id = $1", chatRoom.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
