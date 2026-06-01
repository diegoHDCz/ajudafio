package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}
