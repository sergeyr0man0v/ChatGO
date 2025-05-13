package db

import (
	"chatgo/server/internal/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateUser(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	user := &models.User{
		Username:          "test",
		EncryptedPassword: "password",
		Status:            models.UserStatus("online"),
	}

	rows := sqlmock.NewRows([]string{"id", "username", "encrypted_password", "created_at", "last_login", "status"}).
		AddRow("1", "test", "password", time.Now(), time.Now(), "online")

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.EncryptedPassword, user.Status).
		WillReturnRows(rows)

	ctx := context.Background()
	createdUser, err := repo.CreateUser(ctx, user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "1", createdUser.ID)
	assert.Equal(t, "test", createdUser.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetUserByUsername(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError error
		expectedUser  *models.User
	}{
		{
			name:     "Success",
			username: "test",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "encrypted_password", "created_at", "last_login", "status"}).
					AddRow("1", "test", "password", time.Now(), time.Now(), "online")
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
					WithArgs("test").
					WillReturnRows(rows)
			},
			expectedError: nil,
			expectedUser: &models.User{
				ID:       "1",
				Username: "test",
			},
		},
		{
			name:     "User Not Found",
			username: "test",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
					WithArgs("test").
					WillReturnError(errors.New("user not found"))
			},
			expectedError: errors.New("user not found"),
			expectedUser:  nil,
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
			tc.mockSetup(mock)

			ctx := context.Background()
			user, err := repo.GetUserByUsername(ctx, tc.username)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetUserByID(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "username", "encrypted_password", "created_at", "last_login", "status"}).
		AddRow("1", "test", "password", time.Now(), time.Now(), "online")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	ctx := context.Background()
	user, err := repo.GetUserByID(ctx, "1")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "test", user.Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetAllUsers(t *testing.T) {
	db, mock, err := MockDB(t)
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{"id", "username", "encrypted_password", "created_at", "last_login", "status"}).
		AddRow("1", "test", "password", time.Now(), time.Now(), "online").
		AddRow("2", "test2", "password2", time.Now(), time.Now(), "offline")

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	ctx := context.Background()
	users, err := repo.GetAllUsers(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, "1", users[0].ID)
	assert.Equal(t, "test", users[0].Username)
	assert.Equal(t, "2", users[1].ID)
	assert.Equal(t, "test2", users[1].Username)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
