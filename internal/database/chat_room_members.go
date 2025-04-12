package db

import "chat/internal/models"

// InsertChatroomMember вставляет нового участника чата в базу данных
func (db *Database) InsertChatroomMember(member models.ChatRoomMember) error {
	query := `INSERT INTO chatroom_members (user_id, chat_room_id, joined_at, role) 
        		VALUES ($1, $2, $3, $4)`
	_, err := db.GetDB().Exec(query, member.UserID, member.ChatRoomID, member.JoinedAt, member.Role)
	return err
}

// GetMembersByChatRoom получает участников чата по ID чата
func (db *Database) GetMembersByChatRoom(chatRoomID string) ([]models.ChatRoomMember, error) {
	rows, err := db.GetDB().Query(
		"SELECT user_id, chat_room_id, joined_at, role FROM chatroom_members WHERE chat_room_id = $1",
		chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ChatRoomMember
	for rows.Next() {
		var member models.ChatRoomMember
		if err := rows.Scan(&member.UserID, &member.ChatRoomID, &member.JoinedAt, &member.Role); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}
