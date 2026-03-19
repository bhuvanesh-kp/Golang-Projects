package main

import (
	"contactList/config"
	"contactList/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("Error in loading env variable %v", err)
	}

	PORT := os.Getenv("DB_PORT")

	config.ConnectDB()
	app := gin.Default()
	routes.RegisterRoutes(app)

	app.Run(fmt.Sprintf(":%v", PORT))
}