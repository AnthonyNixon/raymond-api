package auth

import (
	"database/sql"
	"net/http"
	"strings"

	tokens "github.com/AnthonyNixon/raymond-api/cmd/services/token"

	"github.com/AnthonyNixon/raymond-api/cmd/services/database"
	"github.com/AnthonyNixon/raymond-api/cmd/utils/httperr"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(userAuth UserAuth) (token string, error httperr.HttpErr) {
	if userAuth.Password == "" || userAuth.Username == "" {
		return token, httperr.New(http.StatusBadRequest, "Username or Password is empty", "")
	}

	authenticated, err := isAuthed(userAuth.Username, userAuth.Password)
	if err != nil {
		return token, httperr.New(http.StatusBadRequest, "Failed to authenticate user", err.Error())
	}

	if authenticated {
		token, err := tokens.New(userAuth.Username)
		if err != nil {
			return token, err
		}

		return token, nil
	}

	return token, httperr.New(http.StatusUnauthorized, "User Not Authorized", "")
}

func isAuthed(username string, password string) (bool, error) {
	var storedPassword string
	result := database.Connection().QueryRow("select password FROM users where username = ?", username)

	err := result.Scan(&storedPassword)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func AuthenticateRequestAndGetUsername(authHeader string) (username string, error httperr.HttpErr) {
	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return "", httperr.New(http.StatusBadRequest, "bad Authorization header format", "Authorization header is required to be in the format 'Bearer <token>'")
	}

	tokenString := fields[1]

	claims, error := tokens.ParseToken(tokenString)
	if error != nil {
		return username, error
	}

	return claims.Username, nil
}
