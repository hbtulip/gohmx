package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome to bgops,please visit https://github.com/hbtulip",
	})
}
