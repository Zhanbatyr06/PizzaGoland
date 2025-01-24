package users

import (
	"context"

	"github.com/Zhanbatyr06/PizzaGoland/repository/users"

	"testing"

	"github.com/Zhanbatyr06/PizzaGoland/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateUserIntegration(t *testing.T) {
	// 1. Настройка тестовой MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// 2. Используем тестовую базу данных
	db := client.Database("backenddb")
	usersR := users.NewRepo(db, true)

	// 4. Создаём тестовые данные
	inputUser := models.User{
		Nickname: "testuser",
		Password: "securepassword",
	}

	createErr := usersR.Create(&inputUser)
	if createErr != nil {
		t.Fatalf("Failed to create user: %v", createErr)
	}

	var result models.User
	addedUsers, err := usersR.UserFilter("testuser")
	if err != nil {
		t.Fatalf("failed to fetch user from database: %v", err)
	}
	if len(addedUsers) != 1 {
		t.Fatalf("failed to fetch user from database2")
	}
	result = addedUsers[0]

	if result.Nickname != inputUser.Nickname {
		t.Errorf("expected nickname %q, got %q", inputUser.Nickname, result.Nickname)
	}
}
