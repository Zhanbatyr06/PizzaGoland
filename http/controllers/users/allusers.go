package users

import (
	"errors"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func (c Controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := c.repo.GetAll()
	if err != nil {
		utils.Logger.Error(err)
		utils2.SetError(w, errors.New("users not found"))
		return
	}
	utils.Logger.Info("Users fetched successfully")
	utils2.SetJSON(w, users)
}
