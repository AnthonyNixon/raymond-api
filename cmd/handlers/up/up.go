package up

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUpV1(router *gin.Engine) {
	router.GET("/v1/up", GetUp)

}

func GetUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
