package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
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

func EncryptAll(service, username, password, masterPassword string) ([]byte, []byte, []byte) {
	e_service, err := Encrypt(service, masterPassword)
	if err != nil {
		fmt.Printf("Error encrypting service")
		return nil, nil, nil
	}
	e_username, err := Encrypt(username, masterPassword)
	if err != nil {
		fmt.Printf("Error encrypting username")
		return nil, nil, nil
	}
	e_password, err := Encrypt(password, masterPassword)
	if err != nil {
		fmt.Printf("Error encrypting password")
		return nil, nil, nil
	}
	return e_service, e_username, e_password
}

func EncryptInput(service, username, masterPassword string) ([]byte, []byte) {
	e_service, err := Encrypt(service, masterPassword)
	if err != nil {
		fmt.Printf("Error encrypting service")
		return nil, nil
	}
	e_username, err := Encrypt(username, masterPassword)
	if err != nil {
		fmt.Printf("Error encrypting username")
		return nil, nil
	}
	return e_service, e_username
}

func DecryptOutput(service, username []byte, masterPassword string) (string, string, error) {
	dService, err := Decrypt(service, masterPassword)
	if err != nil {
		return "", "", errors.New("Error decrypting service")
	}
	dUsername, err := Decrypt(username, masterPassword)
	if err != nil {
		return "", "", errors.New("Error decrypting username")
	}
	return dService, dUsername, nil
}
