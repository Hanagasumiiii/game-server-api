package handlers

import (
	"encoding/json"
	"game-server-api/internal/user"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func LoginHandler(svc *user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		authenticated, err := svc.AuthenticateUser(req.Username, req.Password)
		if err != nil || !authenticated {
			response := ResponseData{
				Success: false,
				Message: "Invalid username or password",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := ResponseData{
			Success: true,
			Message: "Login successful!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
