package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// converts pass to 32 bytes for AES
func deriveKey(masterPassword string) []byte {
	hash := sha256.Sum256([]byte(masterPassword))
	return hash[:]
}

func Encrypt(plaintext, masterPassword string) ([]byte, error) {
	key := deriveKey(masterPassword)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// nonce is a random number that we only use once
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal encrypts and adds nonce to the beginning
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// Decrypt the data
func Decrypt(ciphertext []byte, masterPassword string) (string, error) {
	key := deriveKey(masterPassword)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Separates nonce from the rest
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("Wrong Master Password or Data corrupted")
	}

	return string(plaintext), nil
}
