package routes

import (
	"github.com/gin-gonic/gin"

	"health_checker/handler"
)

// RegisterRoutes wires all /api/v1 endpoints to their handler functions.
func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	v1.GET("/health", handler.HealthCheck)
	v1.GET("/sites", handler.GetSites)
	v1.POST("/check", handler.RunCheck)
	v1.GET("/history", handler.GetHistory)
	v1.GET("/history/:run_id", handler.GetRunByID)
}
