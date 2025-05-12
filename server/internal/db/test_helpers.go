package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func MockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create mock database: %v", err)
	}

	return db, mock, nil
}

func TestDB(t *testing.T) (*Database, error) {
	t.Helper()
	config := &Config{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "password",
		DBName:   "chat-go-test",
		SSLMode:  "disable",
	}

	db, err := NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clearTables(ctx, db.db)

	return db, nil
}

func clearTables(ctx context.Context, db *sql.DB) {
	tables := []string{"messages", "chat_room_members", "chat_rooms", "users"}

	for _, table := range tables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("Warning: Failed to clear table %s: %v", table, err)
		}
	}
}
