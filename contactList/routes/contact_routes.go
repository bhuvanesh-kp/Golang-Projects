package routes

import (
	"github.com/gin-gonic/gin"
)

func ContactRoutes(app *gin.Engine){
	contacts := app.Group("/contacts")
}