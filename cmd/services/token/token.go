package token

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AnthonyNixon/raymond-api/cmd/utils/httperr"
	jwt "github.com/dgrijalva/jwt-go"
)

var JWT_SIGNING_KEY []byte

const TOKEN_VALID_TIME = 12 * time.Hour

func Initialize() {
	log.Print("Initializing Authentication")
	signingKey := os.Getenv("RAYMOND_JWT_SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("No Signing Key Present.")
	}

	JWT_SIGNING_KEY = []byte(signingKey)
	log.Print("done")
}

func New(username string) (token string, error httperr.HttpErr) {
	var jwtKey = JWT_SIGNING_KEY
	expirationTime := time.Now().Add(TOKEN_VALID_TIME)
	claims := TokenClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString(jwtKey)
	if err != nil {
		return token, httperr.New(http.StatusInternalServerError, "Failed to sign token", err.Error())
	}

	return token, nil

}

func ParseToken(token string) (claims TokenClaims, http_err httperr.HttpErr) {
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SIGNING_KEY, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, httperr.New(http.StatusUnauthorized, "signature invalid", err.Error())
		}

		if !tkn.Valid {
			return claims, httperr.New(http.StatusUnauthorized, "Token invalid", err.Error())
		}

		return claims, httperr.New(http.StatusBadRequest, "Invalid JWT token", err.Error())
	}

	return claims, nil
}
