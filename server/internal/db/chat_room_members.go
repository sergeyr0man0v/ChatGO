package db

import (
	"chatgo/server/internal/models"
	"context"
)

// AddMember вставляет нового участника чата в базу данных
func (r *repository) AddMember(ctx context.Context, member *models.ChatRoomMember) (*models.ChatRoomMember, error) {
	query := `INSERT INTO chat_room_members (user_id, chat_room_id, joined_at, role)
        		VALUES ($1, $2, $3, $4) RETURNING user_id, chat_room_id, joined_at, role`

	err := r.db.QueryRowContext(ctx, query, member.UserID, member.ChatRoomID, member.JoinedAt, member.MemberRole).Scan(
		&member.UserID, &member.ChatRoomID, &member.JoinedAt, &member.MemberRole)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// // GetMembersByUserID получает список вхождений в чаты по ID участника
// func (r *repository) GetMembersByUserID(ctx context.Context, userId string) ([]*models.ChatRoomMember, error) {
// 	rows, err := r.db.QueryContext(ctx,
// 		"SELECT user_id, chat_room_id, joined_at, role FROM chat_room_members WHERE userId = $1",
// 		userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var members []*models.ChatRoomMember
// 	for rows.Next() {
// 		var member models.ChatRoomMember
// 		if err := rows.Scan(&member.UserID, &member.ChatRoomID, &member.JoinedAt, &member.MemberRole); err != nil {
// 			return nil, err
// 		}
// 		members = append(members, &member)
// 	}
// 	return members, nil
// }

// GetMembersByChatRoomID получает участников чата по ID чата
func (r *repository) GetMembersByChatRoomID(ctx context.Context, chatRoomID string) ([]*models.ChatRoomMember, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT user_id, chat_room_id, joined_at, role FROM chat_room_members WHERE chat_room_id = $1",
		chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.ChatRoomMember
	for rows.Next() {
		var member models.ChatRoomMember
		if err := rows.Scan(&member.UserID, &member.ChatRoomID, &member.JoinedAt, &member.MemberRole); err != nil {
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

// // GetMemberByUserAndRoomID получает участника чата по ID чата и ID пользователя
// func (r *repository) GetMemberByUserAndRoomID(ctx context.Context, userID string, chatRoomID string) (*models.ChatRoomMember, error) {
// 	var member models.ChatRoomMember
// 	query := "SELECT user_id, chat_room_id, joined_at, role FROM chat_room_members WHERE user_id = $1 AND chat_room_id = $2"
// 	err := r.db.QueryRowContext(ctx, query, userID, chatRoomID).Scan(&member.UserID, &member.ChatRoomID, &member.JoinedAt, &member.MemberRole)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &member, nil
// }

// UpdateMemberRole обновляет роль участника чата по ID чата и ID пользователя
func (r *repository) UpdateMemberRole(ctx context.Context, member *models.ChatRoomMember) (*models.ChatRoomMember, error) {
	query := "UPDATE chat_room_members SET role = $1 WHERE user_id = $2 AND chat_room_id = $3 RETURNING user_id, chat_room_id, joined_at, role"
	err := r.db.QueryRowContext(ctx, query, member.MemberRole, member.UserID, member.ChatRoomID).Scan(&member.UserID, &member.ChatRoomID, &member.JoinedAt, &member.MemberRole)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// RemoveMember удаляет участника чата по ID чата и ID пользователя
func (r *repository) DeleteMember(ctx context.Context, member *models.ChatRoomMember) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM chat_room_members WHERE user_id = $1 AND chat_room_id = $2", member.UserID, member.ChatRoomID)
	return err
}

// // RemoveMembersByUserID удаляет всех участников чата по ID пользователя
// func (r *repository) RemoveMembersByUserID(ctx context.Context, userID string) error {
// 	_, err := r.db.ExecContext(ctx, "DELETE FROM chat_room_members WHERE user_id = $1", userID)
// 	return err
// }

// RemoveMembersByChatRoomID удаляет всех участников чата по ID чата
func (r *repository) DeleteMembersByChatRoomID(ctx context.Context, chatRoomID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM chat_room_members WHERE chat_room_id = $1", chatRoomID)
	return err
}
