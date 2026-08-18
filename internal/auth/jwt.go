package auth

import (
	"errors"
	"fmt"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL  = 4 * time.Hour
	RefreshTokenTTL = 7 * 24 * time.Hour
	tokenIssuer     = "ajudafio"
)

func GenerateAccessToken(secret []byte, userID, email, name, role string) (string, error) {
	now := time.Now()
	claims := &authdomain.JWTClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(secret []byte, userID string) (signed string, expiresAt time.Time, err error) {
	now := time.Now()
	expiresAt = now.Add(RefreshTokenTTL)
	claims := &authdomain.RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err = token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// HMACKeyfunc returns a jwt.Keyfunc that only accepts HMAC-signed tokens,
// guarding against algorithm-confusion attacks (e.g. "alg: none").
func HMACKeyfunc(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	}
}

func ValidateAccessToken(secret []byte, tokenString string) (*authdomain.JWTClaims, error) {
	claims := &authdomain.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, HMACKeyfunc(secret))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %w", err)
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func ValidateRefreshToken(secret []byte, tokenString string) (*authdomain.RefreshClaims, error) {
	claims := &authdomain.RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, HMACKeyfunc(secret))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("refresh token expired: %w", err)
		}
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}
