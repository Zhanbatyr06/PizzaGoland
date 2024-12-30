package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

// User struct to represent the MongoDB document
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nickname string             `json:"nickname"`
	Password string             `json:"password"`
}

func main() {
	// MongoDB connection setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://Zhanba:UQvsLLt8Tf5lBIvt@godatabase.fzzyv.mongodb.net/go?retryWrites=true&w=majority&appName=GoDatabase"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection = client.Database("go").Collection("users")
	fmt.Println("Connected to MongoDB")

	// Routes
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/user", userHandler)

	// Serve frontend files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler for fetching all users or adding a new user
func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// Fetch all users
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.Background())

		var users []User
		if err = cursor.All(context.Background(), &users); err != nil {
			http.Error(w, fmt.Sprintf(
				"Error decoding users: %v",
				err,
			), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)

	} else if r.Method == http.MethodPost {
		// Add a new user
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result.InsertedID)
	}
}

// Handler for fetching, updating, or deleting a user by ID
func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Fetch user by ID
		var user User
		err := collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		// Update user by ID
		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatedUser})
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("User updated"))

	case http.MethodDelete:
		// Delete user by ID
		_, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("User deleted"))
	}
}
