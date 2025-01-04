package handlers

import (
	"context"
	"encoding/json"
	"github.com/Zhanbatyr06/PizzaGoland/models"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserHandler(client *mongo.Client) http.HandlerFunc {
	collection := client.Database("go").Collection("users")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.URL.Query().Get("id")

		if id == "" {
			utils.Logger.Warn("ID not provided")
			http.Error(w, "ID not provided", http.StatusBadRequest)
			return
		}

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Logger.WithField("id", id).Warn("Invalid ID")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			utils.Logger.WithField("action", "get_user").Info("Fetching user by ID")
			var user models.User
			err := collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
			if err != nil {
				utils.Logger.WithField("id", id).Warn("User not found")
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			utils.Logger.WithField("id", id).Info("User fetched successfully")
			json.NewEncoder(w).Encode(user)

		case http.MethodPut:
			utils.Logger.WithField("action", "update_user").Info("Updating user by ID")
			var updatedUser models.User
			if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
				utils.Logger.WithField("error", err).Warn("Invalid request body")
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			_, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatedUser})
			if err != nil {
				utils.Logger.WithField("error", err).Error("Error updating user")
				http.Error(w, "Error updating user", http.StatusInternalServerError)
				return
			}
			utils.Logger.WithField("id", id).Info("User updated successfully")
			w.Write([]byte("User updated"))

		case http.MethodDelete:
			utils.Logger.WithField("action", "delete_user").Info("Deleting user by ID")
			_, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
			if err != nil {
				utils.Logger.WithField("error", err).Error("Error deleting user")
				http.Error(w, "Error deleting user", http.StatusInternalServerError)
				return
			}
			utils.Logger.WithField("id", id).Info("User deleted successfully")
			w.Write([]byte("User deleted"))

		default:
			utils.Logger.WithField("method", r.Method).Warn("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
