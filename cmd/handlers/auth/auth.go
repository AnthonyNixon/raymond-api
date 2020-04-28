package auth

import (
	"net/http"

	"github.com/AnthonyNixon/raymond-api/cmd/services/auth"
	"github.com/gin-gonic/gin"
)

func AddAuthV1(router *gin.Engine) {
	router.POST("/v1/auth/signin", signIn)
}

func signIn(c *gin.Context) {
	var userSignin auth.UserAuth
	err := c.ShouldBindJSON(&userSignin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind.", "Details": err.Error()})
		return
	}

	token, error := auth.SignIn(userSignin)
	if error != nil {
		c.JSON(error.StatusCode(), error.GetErrorJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
