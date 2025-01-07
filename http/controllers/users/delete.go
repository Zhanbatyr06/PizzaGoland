package users

import (
	"errors"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func (c Controller) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		utils2.SetError(w, errors.New("id is required"))
		return
	}
	deleteErr := c.repo.Delete(id)
	if deleteErr != nil {
		utils2.SetError(w, deleteErr)
	}
	utils.Logger.WithField("id", id).Info("User deleted successfully")
}
