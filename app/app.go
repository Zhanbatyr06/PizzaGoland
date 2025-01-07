package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zhanbatyr06/PizzaGoland/http/controllers"
	"github.com/Zhanbatyr06/PizzaGoland/repository/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	mongoDB *mongo.Database
	usersR  *users.Repo
	httpSrv controllers.Server
}

func NewApp() (*App, error) {
	app := &App{}
	err := app.initMongo()
	if err != nil {
		return nil, err
	}
	app.usersR = users.NewRepo(app.mongoDB)
	app.httpSrv = controllers.New(app.usersR)
	return app, nil
}

func (a *App) Start() {
	// Запуск сервера в горутине
	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	if err := a.httpSrv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited properly")
}

func (a *App) initMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://Zhanba:UQvsLLt8Tf5lBIvt@godatabase.fzzyv.mongodb.net/go?retryWrites=true&w=majority&appName=GoDatabase",
	))
	if err != nil {
		return err
	}
	a.mongoDB = client.Database("go")
	return nil
}
