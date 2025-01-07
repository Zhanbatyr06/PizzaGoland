package users

import (
	"encoding/json"
	"errors"
	"github.com/Zhanbatyr06/PizzaGoland/http/utils"
	utils2 "github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func (c Controller) FilterUsers(w http.ResponseWriter, r *http.Request) {
	// Установка заголовков
	w.Header().Set("Content-Type", "application/json")

	// Извлечение фильтра из строки запроса
	Nickname := r.URL.Query().Get("filter")
	if Nickname == "" {
		utils2.Logger.Warn("Filter parameter is missing in the request")
		utils.SetError(w, errors.New("filter parameter is required"))
		return
	}

	// Логирование полученного фильтра
	utils2.Logger.WithField("filter", Nickname).Info("Processing filter request")

	// Вызов репозитория для фильтрации пользователей
	users, err := c.repo.UserFilter(Nickname)
	if err != nil {
		// Логирование ошибки
		utils2.Logger.WithError(err).WithField("filter", Nickname).Error("Failed to fetch users with the provided filter")

		// Детализированное сообщение клиенту
		errorMessage := map[string]string{
			"error":   "user not found",
			"details": err.Error(), // Подробности ошибки
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorMessage)
		return
	}

	// Проверка, если пользователей не найдено
	if len(users) == 0 {
		utils2.Logger.WithField("filter", Nickname).Info("No users found matching the filter")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "no users found matching the provided filter",
			"filter":  Nickname,
		})
		return
	}

	// Успешное завершение
	utils2.Logger.WithField("filter", Nickname).Info("Users fetched successfully")
	w.WriteHeader(http.StatusOK)
	utils.SetJSON(w, users)
}
