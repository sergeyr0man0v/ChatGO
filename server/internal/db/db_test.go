package db

import (
	"testing"
)

func TestNewDatabase(t *testing.T) {
	testCases := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		// {
		// 	name: "Valid configuration",
		// 	config: &Config{
		// 		Host:     "localhost",
		// 		Port:     "5433",
		// 		User:     "postgres",
		// 		Password: "password",
		// 		DBName:   "chat-go-test",
		// 		SSLMode:  "disable",
		// 	},
		// 	expectError: false,
		// },
		{
			name: "Invalid host",
			config: &Config{
				Host:     "invalid-host",
				Port:     "5433",
				User:     "postgres",
				Password: "password",
				DBName:   "chat-go-test",
				SSLMode:  "disable",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := NewDatabase(tc.config)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if err == nil {
				if db == nil {
					t.Errorf("Expected non-nil database but got nil")
				}
				defer db.Close()

				sqlDB := db.GetDB()
				if sqlDB == nil {
					t.Errorf("GetDB() returned nil")
				}
			}
		})
	}
}
