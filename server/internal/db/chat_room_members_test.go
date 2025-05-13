package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"chatgo/server/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_AddMember(t *testing.T) {
	testCases := []struct {
		name        string
		member      *models.ChatRoomMember
		mockSetup   func(mock sqlmock.Sqlmock, member *models.ChatRoomMember)
		expectError bool
		checkResult func(t *testing.T, member *models.ChatRoomMember, err error)
	}{
		{
			name: "Successfully add member",
			member: &models.ChatRoomMember{
				UserID:     "1",
				ChatRoomID: "1",
				JoinedAt:   time.Now(),
				MemberRole: "member",
			},
			mockSetup: func(mock sqlmock.Sqlmock, member *models.ChatRoomMember) {
				rows := sqlmock.NewRows([]string{"user_id", "chat_room_id", "joined_at", "role"}).
					AddRow(member.UserID, member.ChatRoomID, member.JoinedAt, member.MemberRole)

				mock.ExpectQuery("INSERT INTO chat_room_members").
					WithArgs(member.UserID, member.ChatRoomID, member.JoinedAt, member.MemberRole).
					WillReturnRows(rows)
			},
			expectError: false,
			checkResult: func(t *testing.T, member *models.ChatRoomMember, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, member)
				assert.Equal(t, "1", member.UserID)
				assert.Equal(t, "1", member.ChatRoomID)
			},
		},
		{
			name: "Error when adding member to non-existent chat room",
			member: &models.ChatRoomMember{
				UserID:     "1",
				ChatRoomID: "999",
				MemberRole: "member",
			},
			mockSetup: func(mock sqlmock.Sqlmock, member *models.ChatRoomMember) {
				mock.ExpectQuery("INSERT INTO chat_room_members").
					WithArgs(member.UserID, member.ChatRoomID, sqlmock.AnyArg(), member.MemberRole).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			checkResult: func(t *testing.T, member *models.ChatRoomMember, err error) {
				assert.Error(t, err)
				assert.Nil(t, member)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := MockDB(t)
			if err != nil {
				t.Fatalf("Error creating mock DB: %v", err)
			}
			defer db.Close()

			repo := &repository{db: db}
			tc.mockSetup(mock, tc.member)

			ctx := context.Background()
			addedMember, err := repo.AddMember(ctx, tc.member)

			tc.checkResult(t, addedMember, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetMembersByChatRoomID(t *testing.T) {
	testCases := []struct {
		name        string
		chatRoomID  string
		mockSetup   func(mock sqlmock.Sqlmock, chatRoomID string)
		expectError bool
		checkResult func(t *testing.T, members []*models.ChatRoomMember, err error)
	}{
		{
			name:       "Get members from existing chat room",
			chatRoomID: "1",
			mockSetup: func(mock sqlmock.Sqlmock, chatRoomID string) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{"user_id", "chat_room_id", "joined_at", "role"}).
					AddRow("1", chatRoomID, now, "admin").
					AddRow("2", chatRoomID, now, "member")

				mock.ExpectQuery("SELECT (.+) FROM chat_room_members WHERE chat_room_id = \\$1").
					WithArgs(chatRoomID).
					WillReturnRows(rows)
			},
			expectError: false,
			checkResult: func(t *testing.T, members []*models.ChatRoomMember, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, members)
				assert.Len(t, members, 2)
				assert.Equal(t, "1", members[0].UserID)
				assert.Equal(t, "2", members[1].UserID)
			},
		},
		{
			name:       "Get members from non-existent chat room",
			chatRoomID: "999",
			mockSetup: func(mock sqlmock.Sqlmock, chatRoomID string) {
				mock.ExpectQuery("SELECT (.+) FROM chat_room_members WHERE chat_room_id = \\$1").
					WithArgs(chatRoomID).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "chat_room_id", "joined_at", "role"}))
			},
			expectError: false,
			checkResult: func(t *testing.T, members []*models.ChatRoomMember, err error) {
				assert.NoError(t, err)
				assert.Empty(t, members)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := MockDB(t)
			if err != nil {
				t.Fatalf("Error creating mock DB: %v", err)
			}
			defer db.Close()

			repo := &repository{db: db}
			tc.mockSetup(mock, tc.chatRoomID)

			ctx := context.Background()
			members, err := repo.GetMembersByChatRoomID(ctx, tc.chatRoomID)

			tc.checkResult(t, members, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_UpdateMemberRole(t *testing.T) {
	testCases := []struct {
		name        string
		member      *models.ChatRoomMember
		mockSetup   func(mock sqlmock.Sqlmock, member *models.ChatRoomMember)
		expectError bool
		checkResult func(t *testing.T, member *models.ChatRoomMember, err error)
	}{
		{
			name: "Successfully update member role",
			member: &models.ChatRoomMember{
				UserID:     "1",
				ChatRoomID: "1",
				MemberRole: "moderator",
			},
			mockSetup: func(mock sqlmock.Sqlmock, member *models.ChatRoomMember) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{"user_id", "chat_room_id", "joined_at", "role"}).
					AddRow(member.UserID, member.ChatRoomID, now, member.MemberRole)

				mock.ExpectQuery("UPDATE chat_room_members SET role = \\$1 WHERE user_id = \\$2 AND chat_room_id = \\$3").
					WithArgs(member.MemberRole, member.UserID, member.ChatRoomID).
					WillReturnRows(rows)
			},
			expectError: false,
			checkResult: func(t *testing.T, member *models.ChatRoomMember, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, member)
				assert.Equal(t, "1", member.UserID)
				assert.Equal(t, "moderator", string(member.MemberRole))
			},
		},
		{
			name: "Update role for non-existent member",
			member: &models.ChatRoomMember{
				UserID:     "999",
				ChatRoomID: "1",
				MemberRole: "moderator",
			},
			mockSetup: func(mock sqlmock.Sqlmock, member *models.ChatRoomMember) {
				mock.ExpectQuery("UPDATE chat_room_members SET role = \\$1 WHERE user_id = \\$2 AND chat_room_id = \\$3").
					WithArgs(member.MemberRole, member.UserID, member.ChatRoomID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			checkResult: func(t *testing.T, member *models.ChatRoomMember, err error) {
				assert.Error(t, err)
				assert.Nil(t, member)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := MockDB(t)
			if err != nil {
				t.Fatalf("Error creating mock DB: %v", err)
			}
			defer db.Close()

			repo := &repository{db: db}
			tc.mockSetup(mock, tc.member)

			ctx := context.Background()
			updatedMember, err := repo.UpdateMemberRole(ctx, tc.member)

			tc.checkResult(t, updatedMember, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_DeleteMember(t *testing.T) {
	testCases := []struct {
		name        string
		member      *models.ChatRoomMember
		mockSetup   func(mock sqlmock.Sqlmock, member *models.ChatRoomMember)
		expectError bool
	}{
		{
			name: "Successfully delete member",
			member: &models.ChatRoomMember{
				UserID:     "1",
				ChatRoomID: "1",
			},
			mockSetup: func(mock sqlmock.Sqlmock, member *models.ChatRoomMember) {
				mock.ExpectExec("DELETE FROM chat_room_members WHERE user_id = \\$1 AND chat_room_id = \\$2").
					WithArgs(member.UserID, member.ChatRoomID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name: "Delete non-existent member",
			member: &models.ChatRoomMember{
				UserID:     "999",
				ChatRoomID: "1",
			},
			mockSetup: func(mock sqlmock.Sqlmock, member *models.ChatRoomMember) {
				mock.ExpectExec("DELETE FROM chat_room_members WHERE user_id = \\$1 AND chat_room_id = \\$2").
					WithArgs(member.UserID, member.ChatRoomID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := MockDB(t)
			if err != nil {
				t.Fatalf("Error creating mock DB: %v", err)
			}
			defer db.Close()

			repo := &repository{db: db}
			tc.mockSetup(mock, tc.member)

			ctx := context.Background()
			err = repo.DeleteMember(ctx, tc.member)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
