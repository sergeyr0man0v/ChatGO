package util

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "password123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be the same as the original password")
	}

	err = CheckPassword(password, hashedPassword)
	if err != nil {
		t.Errorf("Expected no error for correct password, got %v", err)
	}

	wrongPassword := "wrong password"
	err = CheckPassword(wrongPassword, hashedPassword)
	if err == nil {
		t.Error("Expected error for incorrect password, got nil")
	}
}
