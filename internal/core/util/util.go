package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateULID() string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}

func HashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hb := hash.Sum(nil)
	return hex.EncodeToString(hb)
}

func EncryptString(str string, secret string) (string, error) {
	key := []byte(secret)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}
	encStr := gcm.Seal(nonce, nonce, []byte(str), nil)

	return hex.EncodeToString(encStr), nil
}

func DecryptString(encStr string, secret string) (string, error) {
	key := []byte(secret)

	cipherText, _ := hex.DecodeString(encStr)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	pt, err := gcm.Open(nil, []byte(nonce), []byte(cipherText), nil)
	if err != nil {
		return "", err
	}
	return string(pt[:]), nil
}
