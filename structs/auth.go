package structs

import "github.com/golang-jwt/jwt"

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	jwt.StandardClaims
}

type Login struct {
	Username string `json:"username" query:"username" form:"username" validate:"required"`
	Password string `json:"password" query:"password" form:"password" validate:"required"`
}
