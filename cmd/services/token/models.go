package token

import jwt "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (t TokenClaims) Valid() error {
	return nil
}
