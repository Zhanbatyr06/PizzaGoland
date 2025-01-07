package users

import (
	"encoding/json"
	"github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"github.com/Zhanbatyr06/PizzaGoland/models"
	"net/http"
)

func (c Controller) createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		utils.SetError(w, decodeErr)
		return
	}
	indesertErr := c.repo.Create(&user)
	if indesertErr != nil {
		utils.SetError(w, indesertErr)
		return
	}
	utils.SetOK(w)
}
