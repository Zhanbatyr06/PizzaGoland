package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"pizzagoland/models"
	"pizzagoland/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UsersHandler(client *mongo.Client) http.HandlerFunc {
	collection := client.Database("go").Collection("users")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			utils.Logger.WithField("action", "fetch_users").Info("Fetching all users")
			cursor, err := collection.Find(context.Background(), bson.D{})
			if err != nil {
				utils.Logger.WithField("error", err).Error("Failed to fetch users")
				http.Error(w, "Error fetching users", http.StatusInternalServerError)
				return
			}
			defer cursor.Close(context.Background())

			var users []models.User
			if err = cursor.All(context.Background(), &users); err != nil {
				utils.Logger.WithField("error", err).Error("Failed to decode users")
				http.Error(w, "Error decoding users", http.StatusInternalServerError)
				return
			}
			utils.Logger.WithField("user_count", len(users)).Info("Users fetched successfully")
			json.NewEncoder(w).Encode(users)

		case http.MethodPost:
			utils.Logger.WithField("action", "add_user").Info("Adding new user")
			var user models.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				utils.Logger.WithField("error", err).Warn("Invalid request body")
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			result, err := collection.InsertOne(context.Background(), user)
			if err != nil {
				utils.Logger.WithField("error", err).Error("Failed to save user")
				http.Error(w, "Error saving user", http.StatusInternalServerError)
				return
			}
			utils.Logger.WithField("user_id", result.InsertedID).Info("User added successfully")
			json.NewEncoder(w).Encode(result.InsertedID)

		default:
			utils.Logger.WithField("method", r.Method).Warn("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
