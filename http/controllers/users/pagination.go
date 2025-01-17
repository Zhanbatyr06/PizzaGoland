package users

import (
	"github.com/Zhanbatyr06/PizzaGoland/http/utils"
	"net/http"
	"strconv"
)

func (c Controller) Pagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем параметры `limit` и `page`
	limitParam := r.URL.Query().Get("limit")
	pageParam := r.URL.Query().Get("page")

	limit := 10 // Значение по умолчанию
	page := 1   // Значение по умолчанию

	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil {
			limit = l
		}
	}
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil {
			page = p
		}
	}
	skip := (page - 1) * limit
	users, err := c.repo.UserPaging(skip, limit)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.SetJSON(w, users)
}
