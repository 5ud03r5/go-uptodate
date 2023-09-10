package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

var TokenAuth *jwtauth.JWTAuth
var RefreshTokenAuth *jwtauth.JWTAuth

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

func GenerateJWTTokens(additionalClaims map[string]interface{}) (string, string, error) {
	claims := make(map[string]interface{})

	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = "uptodate"
	
	for key, value := range additionalClaims {
        claims[key] = value
    }

	_, accessToken, err := TokenAuth.Encode(claims)
	if err != nil {
		fmt.Printf("Error generating access token")
		return "", "", err
	}

	claimsRefresh := make(map[string]interface{})

	claimsRefresh["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	claimsRefresh["iat"] = time.Now().Unix()
	claimsRefresh["iss"] = "uptodate"

	// Only this information is required in refresh token
	for key, value := range additionalClaims {
		claimsRefresh[key] = value   
    }

	_, refreshToken, err := RefreshTokenAuth.Encode(claimsRefresh)
	if err != nil {
		fmt.Printf("Error generating refresh token")
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateJWTAccessToken(sub string, accessType string) (string, error) {
	// Both refresh and access tokens consist of the same claims
	// Difference is in sign key so refresh token cannot be used as an access token
	// This stateless approach allows to keep it simple and fast

	claims := make(map[string]interface{})

	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = sub
	claims["type"] = accessType
	claims["iss"] = "uptodate"

	_, accessToken, err := TokenAuth.Encode(claims)
	if err != nil {
		fmt.Printf("Error generating access token")
		return "", err
	}

	return accessToken, nil
}