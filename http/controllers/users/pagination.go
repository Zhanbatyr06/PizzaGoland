package users

import (
	"errors"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func (c Controller) Pagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	page := 1
	limit := 4

	user, err := c.repo.UserPaging(page, limit)
	if err != nil {
		utils.Logger.Error(err)
		utils2.SetError(w, errors.New("user not found"))
		return
	}
	utils.Logger.WithField("id", page).Info("User fetched successfully")
	utils2.SetJSON(w, user)
}
