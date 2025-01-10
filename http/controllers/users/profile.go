package users

import (
	"net/http"
)

// ProfileHandler обрабатывает страницу профиля пользователя
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./static/profile.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UpdateProfileHandler обрабатывает изменения данных пользователя
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile updated successfully"))
}
