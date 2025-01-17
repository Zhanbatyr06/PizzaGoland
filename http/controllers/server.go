package controllers

import (
	"context"

	"github.com/Zhanbatyr06/PizzaGoland/http/controllers/users"
	users2 "github.com/Zhanbatyr06/PizzaGoland/repository/users"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	usersC     users.Controller
}

func New(usersR *users2.Repo) Server {
	mux := http.NewServeMux()
	usersC := users.New(usersR)
	usersC.Register(mux)
	//staticC := NewStaticController("/static")
	//mux.HandleFunc("/", staticC.ServeHTML)
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	return Server{
		httpServer: &http.Server{
			Addr: ":8080",
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
