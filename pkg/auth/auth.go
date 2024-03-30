package auth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/senyc/jason/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

var KeyLength = 16

func GetNewApiKey() (string, error) {
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

func GetJwtPrivateKey() (*ecdsa.PrivateKey, error) {
	var privateKey *ecdsa.PrivateKey

	pemFilePath := os.Getenv("AUTH_JWT_PEM_PATH")
	if pemFilePath == "" {
		return privateKey, errors.New("No secret file found")
	}
	pemFileBytes, err := os.ReadFile(pemFilePath)
	if err != nil {
		return privateKey, err
	}
	pemBlock, _ := pem.Decode(pemFileBytes)
	privateKey, err = x509.ParseECPrivateKey(pemBlock.Bytes)
	return privateKey, err
}

func GetNewJWT(uuid string) (string, error) {
	var encodedJwt string
	claims := types.JwtClaims{Uuid: uuid}
	j := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	privateKey, err := GetJwtPrivateKey()
	if err != nil {
		return encodedJwt, err
	}
	encodedJwt, err = j.SignedString(privateKey)
	return encodedJwt, err
}
