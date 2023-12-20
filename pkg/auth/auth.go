package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

var KeyLength = 16

func GetApiKey() (string, error) {
	var result string
	b := make([]byte, KeyLength)
	_, err := rand.Read(b)
	if err != nil {
		return result, err
	}

	return base64.RawStdEncoding.EncodeToString(b), err
}

func EncryptApiKey(apiKey string) string {
	h := sha256.New()
	h.Write([]byte(apiKey))
	hashBytes := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashBytes)
}

func IsAuthorized(clearTextPassword string, encryptedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(clearTextPassword))
}

func EncryptPassword(clearText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(clearText), bcrypt.DefaultCost)
	return string(hash), err
}

