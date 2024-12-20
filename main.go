package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// Глобальные переменные для базы данных
var (
	client     *mongo.Client
	collection *mongo.Collection
)

// Структура для пользователя
type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `json:"name"`
	Age  int                `json:"age"`
}

// Структура для обработки входящих данных
type RequestData struct {
	Message string `json:"message"`
}

// Структура для формирования ответа
type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Обработчик для получения всех пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, "Ошибка при чтении данных из MongoDB", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, "Ошибка при обработке данных из MongoDB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Обработчик для добавления нового пользователя
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя в MongoDB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)
}

// Обработчик для получения пользователя по ID
func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Обработчик для обновления пользователя по ID
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var updateData User
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updateData}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Пользователь обновлен")
}

// Обработчик для удаления пользователя по ID
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Пользователь удален")
}

// Обработчик для получения POST и GET запросов с JSON-данными
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var requestData RequestData

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil || requestData.Message == "" {
			sendResponse(w, ResponseData{
				Status:  "fail",
				Message: "Некорректное JSON-сообщение",
			}, http.StatusBadRequest)
			return
		}

		fmt.Println("Получено сообщение:", requestData.Message)

		sendResponse(w, ResponseData{
			Status:  "success",
			Message: "Данные успешно приняты",
		}, http.StatusOK)
	} else if r.Method == http.MethodGet {
		sendResponse(w, ResponseData{
			Status:  "success",
			Message: "Используйте POST для отправки данных",
		}, http.StatusOK)
	} else {
		sendResponse(w, ResponseData{
			Status:  "fail",
			Message: "Метод не поддерживается",
		}, http.StatusMethodNotAllowed)
	}
}

// Функция для отправки JSON-ответа
func sendResponse(w http.ResponseWriter, response ResponseData, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Не удалось подключиться к MongoDB: %v", err)
	}
	fmt.Println("Успешное подключение к MongoDB!")

	// Подключение к коллекции
	database := client.Database("mongo")
	collection = database.Collection("users")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to the PizzaGoland API!")
	})

	// Настройка маршрутов
	http.HandleFunc("/users", getUsersHandler)         // GET /users для получения всех пользователей
	http.HandleFunc("/add_user", addUserHandler)       // POST /add-user для добавления пользователя
	http.HandleFunc("/json", handler)                  // Обработка JSON запросов POST и GET
	http.HandleFunc("/user", getUserByIDHandler)       // GET /user?id= для получения пользователя по ID
	http.HandleFunc("/update-user", updateUserHandler) // PUT /update-user?id= для обновления пользователя
	http.HandleFunc("/delete-user", deleteUserHandler) // DELETE /delete-user?id= для удаления пользователя

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
