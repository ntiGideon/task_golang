package userModel

import "github.com/dgrijalva/jwt-go"

type RegisterUserModel struct {
	Name     string `json:"name" validate:"required,min=5,max=32"`
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginUserModel struct {
	EmailOrUsername string `json:"email_or_username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type VerifyAccountModel struct {
	EmailToken string `json:"email_token" validate:"required"`
}

type JWTPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Id       int    `json:"id"`
}

type JWTClaim struct {
	Id         int         `json:"id"`
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	CustomData interface{} `json:"customData"`
	jwt.StandardClaims
}

type BeginPasswordReset struct {
	Email string `json:"email" validate:"required,email"`
}
