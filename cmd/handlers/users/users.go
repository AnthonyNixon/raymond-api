package users

import (
	"net/http"

	"github.com/AnthonyNixon/raymond-api/cmd/services/users"
	"github.com/gin-gonic/gin"
)

func AddUsersV1(router *gin.Engine) {
	router.POST("/v1/users", newUser)
}

func newUser(c *gin.Context) {
	var newUser users.User
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind.", "Details": err.Error()})
		return
	}

	error := users.NewUser(newUser)
	if error != nil {
		c.JSON(error.StatusCode(), error.GetErrorJSON())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status": "Created"})
}
