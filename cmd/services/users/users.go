package users

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/AnthonyNixon/raymond-api/cmd/services/database"
	"github.com/AnthonyNixon/raymond-api/cmd/utils/httperr"
)

func (user User) IsUnique() (unique bool, error httperr.HttpErr) {
	var count int
	err := database.Connection().QueryRow("SELECT COUNT(*) FROM users where username = ? OR email = ?", user.Username, user.Email).Scan(&count)
	if err != nil {
		return false, httperr.New(http.StatusInternalServerError, "Failed to query if username is unique", err.Error())
	}

	return count == 0, nil
}

func NewUser(newUser User) (error httperr.HttpErr) {
	if newUser.Email == "" || newUser.Username == "" {
		return httperr.New(http.StatusBadRequest, "Must include email and username", "")
	}

	unique, error := newUser.IsUnique()
	if error != nil {
		return error
	}

	if !unique {
		return httperr.New(http.StatusConflict, "Username or Email is already taken", "")
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, "Error hashing password", err.Error())
	}

	stmt, err := database.Connection().Prepare("insert into users (username, email, password, first_name, last_name) values(?,?,?,?,?);")
	if err != nil {
		return httperr.New(http.StatusInternalServerError, "Failed to prepare DB statement", err.Error())
	}

	_, err = stmt.Exec(newUser.Username, newUser.Email, hashedPassword, newUser.FirstName, newUser.LastName)
	if err != nil {
		return httperr.New(http.StatusInternalServerError, "Failed to execute DB Query", err.Error())
	}

	return nil
}
