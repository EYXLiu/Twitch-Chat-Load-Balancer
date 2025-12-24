// health handlers
// 	health endpoint 


package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(router *gin.Engine) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
