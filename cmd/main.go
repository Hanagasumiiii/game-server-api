package main

import (
	"fmt"
	"game-server-api/internal/config"
	"game-server-api/internal/database"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	path, err := filepath.Abs("config.yml")
	if err != nil {
		log.Fatalf("Error getting absolute path: %s", err)
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("Config file not found at path: %s", path)
	} else if err != nil {
		log.Fatalf("Error accessing config file: %s", err)
	}

	log.Printf("Loading configuration from path: %s", path)

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err)
	}

	db := database.NewConnection(*cfg)

	//userService := user.NewService(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running!")
	})
	http.HandleFunc("/testdb", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "DB connection failed: "+err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, "DB connection successful!")
		}
	})

	log.Println("Starting server on port", cfg.Postgres.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Postgres.Port), nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}
