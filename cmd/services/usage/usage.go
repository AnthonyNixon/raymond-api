package usage

import (
	"database/sql"
	"net/http"

	"github.com/AnthonyNixon/raymond-api/cmd/services/database"
	"github.com/AnthonyNixon/raymond-api/cmd/utils/httperr"
)

func GetTokensForUsername(username string) (tokens int, error httperr.HttpErr) {
	err := database.Connection().QueryRow("SELECT tokens FROM users where username = ?", username).Scan(&tokens)
	if err != nil {
		if err == sql.ErrNoRows {
			return tokens, httperr.New(http.StatusNotFound, "DB entry for user not found", "")
		}
		return tokens, httperr.New(http.StatusInternalServerError, "Failed to query user", err.Error())
	}

	return tokens, nil
}

func UseTokenForUsername(username string) (error httperr.HttpErr) {
	stmt, err := database.Connection().Prepare("UPDATE users SET tokens = tokens - 1 WHERE username = ?")
	if err != nil {
		return httperr.New(http.StatusInternalServerError, "Failed to prepare DB statement", err.Error())
	}

	_, err = stmt.Exec(username)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, "Failed to execute DB Query", err.Error())
	}

	return nil
}
