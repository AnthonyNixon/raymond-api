package httperr

import "github.com/gin-gonic/gin"

// New returns an error that formats as the given text.
func New(status int, description string, errorMessage string) HttpErr {
	return &error{status, description, errorMessage}
}

type error struct {
	status       int
	description  string
	errorMessage string
}

func (e *error) Error() string {
	return e.description
}

func (e *error) Description() string {
	return e.description
}

func (e *error) StatusCode() int {
	return e.status
}

func (e *error) ErrorMessage() string {
	return e.errorMessage
}

func (e *error) GetErrorJSON() gin.H {
	return gin.H{"Status": e.StatusCode(), "Error": e.Description(), "Details": e.errorMessage}
}

type HttpErr interface {
	StatusCode() int
	Description() string
	ErrorMessage() string
	GetErrorJSON() gin.H
}
