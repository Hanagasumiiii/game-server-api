package handlers

import (
	"encoding/json"
	"game-server-api/internal/inventory"
	"log"
	"net/http"
	"strconv"
)

type InventoryRequest struct {
	UserID   int `json:"user_id"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

type InventoryResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Items   []inventory.Item `json:"items,omitempty"`
}

func AddItemHandler(svc *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InventoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = svc.AddItem(req.UserID, req.ItemID, req.Quantity)
		if err != nil {
			response := InventoryResponse{
				Success: false,
				Message: "Failed to add item: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := InventoryResponse{
			Success: true,
			Message: "Item added successfully!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func RemoveItemHandler(svc *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InventoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = svc.RemoveItem(req.UserID, req.ItemID, req.Quantity)
		if err != nil {
			response := InventoryResponse{
				Success: false,
				Message: "Failed to remove item: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := InventoryResponse{
			Success: true,
			Message: "Item removed successfully!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetUserItemsHandler(svc *inventory.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		items, err := svc.GetUserItems(userID)
		if err != nil {
			response := InventoryResponse{
				Success: false,
				Message: "Failed to get user items: " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := InventoryResponse{
			Success: true,
			Message: "Items retrieved successfully!",
			Items:   items,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
