package users

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Zhanbatyr06/PizzaGoland/models"
	"github.com/Zhanbatyr06/PizzaGoland/repository/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	requestData := models.User{
		Nickname: "John Doe123",
		Password: "password123",
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database("backenddb")
	usersR := users.NewRepo(db, true)
	controller := Controller{repo: usersR}
	requestBody, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/add_user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.createUser)
	handler.ServeHTTP(rr, req)

	// Check for expected status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Проверяем тело ответа (опционально)
	expectedBody := "User created successfully"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}
