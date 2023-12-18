package auth

import (
	"crypto/rand"
	"encoding/hex"
)
var KeyLength = 32

func GetApiKey() (error, string) {
	var result string
	b := make([]byte, KeyLength)
	_, err := rand.Read(b)
	if err != nil {
		return err, result
	}

	return err, hex.EncodeToString(b)
}
