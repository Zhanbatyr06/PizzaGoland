package users

import (
	"github.com/Zhanbatyr06/PizzaGoland/repository/users"
	"net/http"
)

type Controller struct {
	repo *users.Repo
}

func New(repo *users.Repo) *Controller {
	return &Controller{repo: repo}
}

func (c Controller) Register(mux *http.ServeMux) {
	mux.HandleFunc("/users/get", c.getUser)
	mux.HandleFunc("/add_user", c.createUser)
	mux.HandleFunc("/delete_user", c.deleteUser)
	mux.HandleFunc("/update_user", c.UpdateUser)
	mux.HandleFunc("/filter_user", c.FilterUsers)
	mux.HandleFunc("/all_users", c.GetAllUsers)
	mux.HandleFunc("/sort_user", c.SortingUser)
}
