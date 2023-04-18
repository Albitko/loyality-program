package entities

import (
	"github.com/golang-jwt/jwt/v4"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID       string
	Login    string
	Password string
}

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	jwt.StandardClaims
}
