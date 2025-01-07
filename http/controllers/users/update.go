package users

import (
	"encoding/json"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"github.com/Zhanbatyr06/PizzaGoland/models"
	
	"net/http"
)

func (c Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		utils2.SetError(w, decodeErr)
		return
	}
	updateErr := c.repo.Update(&user)
	if updateErr != nil {

		utils2.SetError(w, updateErr)
		return
	}
}
