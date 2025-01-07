package users

import (
	"errors"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func (c Controller) getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	if id == "" {
		utils2.SetError(w, errors.New("id is required"))
		return
	}
	user, err := c.repo.Get(id)
	if err != nil {
		utils.Logger.Error(err)
		utils2.SetError(w, errors.New("user not found"))
		return
	}
	utils.Logger.WithField("id", id).Info("User fetched successfully")
	utils2.SetJSON(w, user)
}
