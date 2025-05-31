package main

import (
	"go_link_reducer/cmd/api"
	"go_link_reducer/cmd/db"
	"go_link_reducer/config"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	log.Println("getting environment")
	if err := godotenv.Load(); err != nil {
		log.Panic("error getting environment", err)
	}

	cfg := config.InitConfig()

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Panic("failed connect to database", err)
	}

	pingDB(db)

	router := api.NewRoutes(db)

	port := config.GetEnv("PORT", "8080")

	addr := "0.0.0.0:" + port

	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Println("Starting the server on port :8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Server failed to start: ", err)
	}
}

func pingDB(DB *gorm.DB) {
	psqlDB, err := DB.DB()
	if err != nil {
		log.Panic("failed to get db connection", err)

	}

	if err := psqlDB.Ping(); err != nil {
		log.Panic("failed to ping db", err)
	}
}
