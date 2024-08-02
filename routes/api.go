package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fullstack-starter/database"
	"go-fullstack-starter/internal"
	"net/http"
	"os"
)

func RegisterAPIroutes(multiplexer *http.ServeMux) {
	multiplexer.HandleFunc("/api/health", HealthCheckHandler)
	multiplexer.HandleFunc("/api/auth/login", UserLoginHandler)
	multiplexer.HandleFunc("/api/auth/register", UserRegistrationHandler)
}

var db *sql.DB

func GetDBConnection() *sql.DB {
	if db != nil {
		return db
	}

	dbConn, err := database.NewConnection()
	if err != nil {
		fmt.Println("Error connecting to database")
		os.Exit(1)
	}

	return dbConn
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"version": "0.0.1",
	})
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	db = GetDBConnection()
	internal.UserLogin(w, r, db)
}

func UserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	db = GetDBConnection()
	internal.UserRegistration(w, r, db)
}
