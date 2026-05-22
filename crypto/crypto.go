package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// convierte pass en clave de 32 bytes para AES
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
	// nonce es un número aleatorio que usamos una sola vez
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal encripta y agrega el nonce al inicio
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// Decrypt desencripta los datos
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
		return "", errors.New("ciphertext demasiado corto")
	}

	// Separamos el nonce del resto del mensaje encriptado
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("master password incorrecta o datos corruptos")
	}

	return string(plaintext), nil
}
