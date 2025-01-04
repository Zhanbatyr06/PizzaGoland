package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pizzagoland/handlers"
	"pizzagoland/middlewares"
)

func main() {
	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://Zhanba:UQvsLLt8Tf5lBIvt@godatabase.fzzyv.mongodb.net/go?retryWrites=true&w=majority&appName=GoDatabase"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Обслуживание статических файлов
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// API-обработчики
	http.Handle("/api/users", middlewares.RateLimiter(http.HandlerFunc(handlers.UsersHandler(client))))
	http.Handle("/api/user", middlewares.RateLimiter(http.HandlerFunc(handlers.UserHandler(client))))

	// Настройка сервера
	srv := &http.Server{
		Addr: ":8080",
	}

	// Запуск сервера в горутине
	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Захват системных сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited properly")
}
