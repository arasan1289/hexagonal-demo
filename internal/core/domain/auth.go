package domain

import "github.com/golang-jwt/jwt/v5"

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type UserClaims struct {
	Role string `json:"role"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}
