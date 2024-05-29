package main

import (
	"fmt"
	"game-server-api/internal/config"
	"game-server-api/internal/database"
	"game-server-api/internal/handlers"
	"game-server-api/internal/inventory"
	"game-server-api/internal/user"

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
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}

	userService := user.NewService(db)
	inventoryService := inventory.NewService(db)

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

	http.HandleFunc("/api/login", handlers.LoginHandler(userService))
	http.HandleFunc("/api/register", handlers.RegisterHandler(userService))
	http.HandleFunc("/api/additem", handlers.AddItemHandler(inventoryService))
	http.HandleFunc("/api/removeitem", handlers.RemoveItemHandler(inventoryService))
	http.HandleFunc("/api/getitems", handlers.GetUserItemsHandler(inventoryService))

	log.Println("Starting server on port", cfg.Postgres.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Postgres.Port), nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
