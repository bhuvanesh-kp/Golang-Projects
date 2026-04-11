package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *gin.Engine){
	api := app.Group("/api")
	fmt.Println(api)
}