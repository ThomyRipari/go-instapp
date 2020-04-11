package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	FirstName string `json:"firstname"`
	SurName   string `json:"surname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ctxKey string

const ContextTokenClaims ctxKey = "token"
