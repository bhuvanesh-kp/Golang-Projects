package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ContactRoutes(app *gin.Engine){
	contacts := app.Group("/contacts")
	fmt.Println(contacts)
}