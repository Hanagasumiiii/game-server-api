package handlers

import (
	"encoding/json"
	"game-server-api/internal/user"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = svc.CreateUser(req.Username, req.Email, req.Password)
		if err != nil {
			response := ResponseData{
				Success: false,
				Message: "Failed to create user: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := ResponseData{
			Success: true,
			Message: "User created successfully!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
