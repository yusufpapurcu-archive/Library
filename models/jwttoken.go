package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	UserId string
	jwt.StandardClaims
}
