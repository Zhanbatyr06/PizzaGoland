package users

import (
	"errors"
	"github.com/Zhanbatyr06/PizzaGoland/http/utils"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func (c Controller) SortingUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	SortType := r.URL.Query().Get("sort")
	users, err := c.repo.UserSorting(SortType)
	if err != nil {
		utils2.Logger.WithError(err).WithField("sort", SortType).Error("Failed to fetch users with the provided sorting")
		return
	}
	if SortType == "" {
		utils2.Logger.Warn("Sorting parameter is missing in the request")
		utils.SetError(w, errors.New("Sorting parameter is required"))
		return
	}
	utils2.Logger.WithField("sort", SortType).Info("Processing sorting request")
	w.WriteHeader(http.StatusOK)
	utils.SetJSON(w, users)
}
