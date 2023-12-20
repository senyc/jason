package auth

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

var KeyLength = 32

func GetApiKey() (string, error) {
	var result string
	b := make([]byte, KeyLength)
	_, err := rand.Read(b)
	if err != nil {
		return result, err
	}

	return hex.EncodeToString(b), err
}

func IsAuthorized(clearTextPassword string, encryptedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(clearTextPassword))
}

func EncryptPassword(clearText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(clearText), bcrypt.DefaultCost)
	return string(hash), err
}
