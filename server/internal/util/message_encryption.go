package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"log"
)

// EncryptMessage encrypts a string using AES-256 encryption
func EncryptMessage(message string, key []byte) (string, error) {
	if len(key) != 32 {
		log.Print("Key is not 32 bytes long")
		return message, nil
	}
	plaintext := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	// Convert to base64 for easy transmission
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptMessage decrypts an encrypted string using AES-256 decryption
func DecryptMessage(encryptedMessage string, key []byte) (string, error) {
	if len(key) != 32 {
		log.Print("Key is not 32 bytes long")
		return encryptedMessage, nil
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateKey generates a random 32-byte key for AES-256
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 requires 32-byte key
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
