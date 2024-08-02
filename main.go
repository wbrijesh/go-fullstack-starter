package main

import (
	"fmt"
	"go-fullstack-starter/database"
	"go-fullstack-starter/middleware"
	"go-fullstack-starter/routes"
	"go-fullstack-starter/utils"
	"log"
	"net/http"
	"os"
)

func init() {
	config := utils.ReadConfig()
	if config.Port == 0 || config.JwtKey == "" || config.RateLimitRequestsPerSecond == 0 {
		fmt.Println("Invalid configuration")
		os.Exit(1)
	}
	fmt.Println("Configuration loaded successfully")

	db, err := database.NewConnection()
	if err != nil {
		fmt.Println("Error connecting to database")
		os.Exit(1)
	}
	fmt.Println("Database connection successful")

	database.InitialiseDatabase(db)
}

func main() {
	config := utils.ReadConfig()
	multiplexer := http.NewServeMux()

	routes.RegisterAPIroutes(multiplexer)

	multiplexerWithMiddleware := middleware.RequestIDMiddleware(multiplexer)

	fmt.Printf("\nServer listening on port %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), multiplexerWithMiddleware))
}
