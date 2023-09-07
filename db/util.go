package db

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// Id is in form of:
// Hash@name%version
// name is lowercase and contains - instead of spaces
// version contains _ instead of dots
// Hashing is for randomness

func ConvertToId(name string, version string) string {
	concatenated := name + "%" + version
	resultNoHyphens := strings.ReplaceAll(concatenated, " ", "-")
	resultNoDots := strings.ReplaceAll(resultNoHyphens, ".", "_")
	result := strings.ToLower(resultNoDots)

	hasher := md5.New()
	hasher.Write([]byte(result))
	hash := hex.EncodeToString(hasher.Sum(nil))
	
	id := hash + "@" + result 
	return id
}

func ConvertFromId(id string) (string, string) {
	parts := strings.Split(id, "@")

	// Ensure that we have at least two parts
	if len(parts) < 2 {
		fmt.Printf("Invalid id: %s", id)
		return "", ""
	}
	
	firstPart := parts[1]
	subparts := strings.Split(firstPart, "%")

	// Ensure we have two values in subparts
	if len(subparts) < 2 {
		fmt.Printf("Issue parsing id: %s", id)
		return "", ""
	}
	name := subparts[0]
	version := subparts[1]
	name = strings.ReplaceAll(name, "-", " ")
	version = strings.ReplaceAll(version, "_", ".")
	return name, version	
}










