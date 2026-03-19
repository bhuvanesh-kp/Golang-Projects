package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB(){
	var err error
	fmt.Println("Connecting to postgres database ...")

	postgres_connectionURL := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", 
				os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			 	os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	
	db, err = gorm.Open(postgres.Open(postgres_connectionURL), &gorm.Config{})
	
	if err != nil{
		log.Fatalf("Error connecting db: %v", err.Error())
	}

	fmt.Println("Connection the DB sucessfully")
}