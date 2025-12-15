package handlers

import (
	"net/http"
	"tc/internal/metrics"
	"tc/internal/worker"

	"github.com/gin-gonic/gin"
)

func AdminStatsHandler(router *gin.Engine, counter *metrics.Counter, workers *[]*worker.Worker, msgQueue chan string) {
	router.GET("/admin/stats", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"dropped_messages": counter.Get(),
			"active_workers":   len(*workers),
			"queue_length":     len(msgQueue),
			"queue_capacity":   cap(msgQueue),
		})
	})
}
