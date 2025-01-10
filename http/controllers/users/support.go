package users

import (
	utils2 "github.com/Zhanbatyr06/PizzaGoland/utils"
	"log"
	"net/http"
)

// SupportHandler обрабатывает отправку сообщений в поддержку
func SupportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Handling POST request at /support")

	email := r.FormValue("email")
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	if email == "" || subject == "" || message == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	log.Printf("Email: %s, Subject: %s, Message: %s", email, subject, message)

	err := utils2.SendEmail(email, subject, message, "")
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
