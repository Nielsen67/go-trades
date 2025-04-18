package main

import (
	"go-trades/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading Environment")
	}

	//SETUP CONNECTION DB
	db := config.ConnectDatabase()

	//SETUP MIGRATION
	config.Migrate(db)

	//SETUP ROUTE
	r := routeInit(db)
	r.Run(":3000")
}
