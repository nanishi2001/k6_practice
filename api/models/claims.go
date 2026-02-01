package models

import "github.com/golang-jwt/jwt/v5"

// Claims はJWTトークンに含まれるユーザー情報
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
