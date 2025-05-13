package db

import (
	"chatgo/server/internal/models"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateChatRoom(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	chatRoom := &models.ChatRoom{
		Name:      "Test Room",
		Type:      "group",
		CreatorID: "1",
	}

	mock.ExpectBegin()

	roomRows := sqlmock.NewRows([]string{"id", "name", "type", "creator_id", "created_at"}).
		AddRow("1", "Test Room", "group", "1", time.Now())

	mock.ExpectQuery("INSERT INTO chat_rooms").
		WithArgs(chatRoom.Name, chatRoom.Type, chatRoom.CreatorID).
		WillReturnRows(roomRows)

	memberRows := sqlmock.NewRows([]string{"user_id", "chat_room_id", "member_role", "joined_at"}).
		AddRow("1", "1", "admin", time.Now())

	mock.ExpectQuery("INSERT INTO chat_room_members").
		WithArgs("1", "1", "admin").
		WillReturnRows(memberRows)

	mock.ExpectCommit()

	ctx := context.Background()
	createdRoom, err := repo.CreateChatRoom(ctx, chatRoom)

	assert.NoError(t, err)
	assert.NotNil(t, createdRoom)
	assert.Equal(t, "1", createdRoom.ID)
	assert.Equal(t, "Test Room", createdRoom.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetChatRoomByID(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "name", "type", "creator_id", "created_at"}).
		AddRow("1", "Test Room", "group", "1", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM chat_rooms WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	ctx := context.Background()
	room, err := repo.GetChatRoomByID(ctx, "1")

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, "1", room.ID)
	assert.Equal(t, "Test Room", room.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_UpdateChatRoom(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	chatRoom := &models.ChatRoom{
		ID:        "1",
		Name:      "Updated Room",
		Type:      "group",
		CreatorID: "1",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "type", "creator_id", "created_at"}).
		AddRow("1", "Updated Room", "group", "1", time.Now())

	mock.ExpectQuery("UPDATE chat_rooms SET name = \\$1, type = \\$2 WHERE id = \\$3").
		WithArgs(chatRoom.Name, chatRoom.Type, chatRoom.ID).
		WillReturnRows(rows)

	ctx := context.Background()
	updatedRoom, err := repo.UpdateChatRoom(ctx, chatRoom)

	assert.NoError(t, err)
	assert.NotNil(t, updatedRoom)
	assert.Equal(t, "1", updatedRoom.ID)
	assert.Equal(t, "Updated Room", updatedRoom.Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_DeleteChatRoom(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	chatRoom := &models.ChatRoom{
		ID: "1",
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM chat_room_members WHERE chat_room_id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM chat_rooms WHERE id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	ctx := context.Background()
	err = repo.DeleteChatRoom(ctx, chatRoom)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
