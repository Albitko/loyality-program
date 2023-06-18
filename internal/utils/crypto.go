package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

const secretLength = 20

func HexHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateSecret() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	secret := make([]rune, secretLength)
	for i := range secret {
		secret[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(secret)
}
