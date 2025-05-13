package db

import (
	"context"
	"testing"
	"time"

	"chatgo/server/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateMessage(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	// replyID := sql.NullInt64{Int64: 0, Valid: false}
	message := &models.Message{
		SenderID:         "1",
		ChatRoomID:       "1",
		EncryptedContent: "Test message content",
		// ReplyToMessageID: replyID,
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "chat_room_id", "encrypted_content" /*"reply_to_message_id",*/, "created_at", "updated_at", "is_edited"}).
		AddRow("1", "1", "1", "Test message content" /*replyID,*/, time.Now(), time.Now(), false)

	mock.ExpectQuery("INSERT INTO messages").
		WithArgs(message.SenderID, message.ChatRoomID, message.EncryptedContent /*message.ReplyToMessageID*/).
		WillReturnRows(rows)

	ctx := context.Background()
	createdMessage, err := repo.CreateMessage(ctx, message)

	assert.NoError(t, err)
	assert.NotNil(t, createdMessage)
	assert.Equal(t, "1", createdMessage.ID)
	assert.Equal(t, "Test message content", createdMessage.EncryptedContent)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetMessagesByChatRoomID(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "chat_room_id", "encrypted_content", "created_at", "updated_at", "is_edited"}).
		AddRow("1", "1", "1", "Message 1", time.Now(), time.Now(), false).
		AddRow("2", "2", "1", "Message 2", time.Now(), time.Now(), false)

	mock.ExpectQuery("SELECT (.+) FROM messages WHERE chat_room_id = \\$1 ORDER BY created_at ASC LIMIT \\$2").
		WithArgs("1", 10).
		WillReturnRows(rows)

	ctx := context.Background()
	messages, err := repo.GetMessagesByChatRoomID(ctx, "1", 10)

	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Len(t, messages, 2)
	assert.Equal(t, "1", messages[0].ID)
	assert.Equal(t, "2", messages[1].ID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetMessageByID(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "chat_room_id", "encrypted_content", "created_at", "updated_at", "is_edited"}).
		AddRow("1", "1", "1", "Test message", time.Now(), time.Now(), false)

	mock.ExpectQuery("SELECT (.+) FROM messages WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	ctx := context.Background()
	message, err := repo.GetMessageByID(ctx, "1")

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "1", message.ID)
	assert.Equal(t, "Test message", message.EncryptedContent)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
