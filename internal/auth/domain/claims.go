package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Exp               int64              `json:"exp"`
	Iat               int64              `json:"iat"`
	AuthTime          int64              `json:"auth_time"`
	Jti               string             `json:"jti"`
	Iss               string             `json:"iss"`
	Aud               string             `json:"aud"`
	Sub               string             `json:"sub"`
	Typ               string             `json:"typ"`
	Azp               string             `json:"azp"`
	Sid               string             `json:"sid"`
	Acr               string             `json:"acr"`
	AllowedOrigins    []string           `json:"allowed-origins"`
	RealmAccess       RealmAccess        `json:"realm_access"`
	ResourceAccess    map[string]RoleSet `json:"resource_access"`
	Scope             string             `json:"scope"`
	EmailVerified     bool               `json:"email_verified"`
	Name              string             `json:"name"`
	PreferredUsername string             `json:"preferred_username"`
	GivenName         string             `json:"given_name"`
	FamilyName        string             `json:"family_name"`
	Email             string             `json:"email"`
	jwt.RegisteredClaims
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type RoleSet struct {
	Roles []string `json:"roles"`
}
