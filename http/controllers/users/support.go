package users

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func SupportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	file, header, err := r.FormFile("attachment")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Сохранение файла
	if file != nil {
		defer file.Close()
		out, err := os.Create("./uploads/" + header.Filename)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		io.Copy(out, file)
	}

	// Логика отправки письма или записи в базу данных
	fmt.Printf("Received support request: Email: %s, Subject: %s, Message: %s\n", email, subject, message)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Support request submitted"))
}
