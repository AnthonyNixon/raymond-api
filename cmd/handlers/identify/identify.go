package identify

import (
	"fmt"
	"net/http"

	"github.com/AnthonyNixon/raymond-api/cmd/services/audd"

	"github.com/AnthonyNixon/raymond-api/cmd/services/auth"
	"github.com/AnthonyNixon/raymond-api/cmd/services/usage"
	"github.com/AnthonyNixon/raymond-api/cmd/utils/httperr"
	"github.com/gin-gonic/gin"
)

func AddIdentifyV1(router *gin.Engine) {
	router.POST("/v1/identify", identify)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

}

func identify(c *gin.Context) {
	username, err := auth.AuthenticateRequestAndGetUsername(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(err.StatusCode(), err.GetErrorJSON())
		return
	}

	tokensAvailable, err := usage.GetTokensForUsername(username)
	if err != nil {
		c.JSON(err.StatusCode(), err.GetErrorJSON())
		return
	}

	if tokensAvailable == 0 {
		usageError := httperr.New(http.StatusPaymentRequired, "Tokens required", fmt.Sprintf("User %s has no tokens available to use.", username))
		c.JSON(http.StatusPaymentRequired, usageError.GetErrorJSON())
	}

	fileHeader, error := c.FormFile("file")
	if error != nil {
		fileError := httperr.New(http.StatusBadRequest, "Failed to get file from request", error.Error())
		c.JSON(fileError.StatusCode(), fileError.GetErrorJSON())
		return
	}

	//filename := filepath.Base(file.Filename)
	file, error := fileHeader.Open()
	if error != nil {
		fileError := httperr.New(http.StatusBadRequest, "Failed to open file", error.Error())
		c.JSON(fileError.StatusCode(), fileError.GetErrorJSON())
		return
	}

	// User has usage remaining
	title, artist, err := audd.IdentifyFile(file)
	if err != nil {
		c.JSON(err.StatusCode(), err.GetErrorJSON())
		return
	}

	err = usage.UseTokenForUsername(username)
	if err != nil {
		c.JSON(err.StatusCode(), err.GetErrorJSON())
		return
	}

	c.JSON(http.StatusOK, gin.H{"Title": title, "Artist": artist})
}
