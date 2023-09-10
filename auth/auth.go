package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

var TokenAuth *jwtauth.JWTAuth

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}

func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
		fmt.Print("Error while hashing password")
        return "", err
    }
    return string(hashedPassword), nil
}

func VerifyPassword(inputPassword string, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
    return err == nil
}

func GeneratePassword(length int) (string, error) {
	// Generate random data
	randomData := make([]byte, length)
	_, err := rand.Read(randomData)
	if err != nil {
		fmt.Print("Error while generating password")
		return "", err
	}

	// Calculate the SHA-256 hash of the random data
	hash := sha256.Sum256(randomData)

	// Encode the hash in hexadecimal
	password := hex.EncodeToString(hash[:])

	return password, nil
}