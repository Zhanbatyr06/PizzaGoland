package controllers

import (
	"context"
	"net/http"

	"github.com/Zhanbatyr06/PizzaGoland/http/controllers/users"
	users2 "github.com/Zhanbatyr06/PizzaGoland/repository/users"
)

type Server struct {
	httpServer *http.Server
	usersC     users.Controller
}

func New(usersR *users2.Repo) Server {
	// Создаем новый маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем контроллеры пользователей
	usersC := users.New(usersR)
	usersC.Register(mux)

	// Регистрируем обработчик для статических файлов
	staticPath := "D:/AITUAssignments/PizzaGoland/static"
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	// Регистрируем обработчик для главной страницы
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticPath+"/index.html")
	})

	// Создаем и возвращаем сервер
	return Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux, // Назначаем маршрутизатор обработчиком
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
