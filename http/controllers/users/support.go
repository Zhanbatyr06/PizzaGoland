package users

import (
	"fmt"
	"github.com/go-gomail/gomail"
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
	var filePath string
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Сохранение файла
	if file != nil {
		defer file.Close()
		filePath = "./uploads/" + header.Filename
		if err := saveFile(file, filePath); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
	}

	// Логика отправки письма
	if err := sendEmailWithAttachment(email, subject, message, filePath); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Support request submitted"))
}

// saveFile сохраняет файл на диск
func saveFile(file io.Reader, path string) error {
	// Убедитесь, что директория существует
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create uploads directory: %w", err)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// sendEmailWithAttachment отправляет письмо с вложением
func sendEmailWithAttachment(fromEmail, subject, body, attachmentPath string) error {
	// Настройки SMTP сервера
	const (
		smtpHost = "smtp.gmail.com"
		smtpPort = 587
		username = "zhanbatyrmolkryt@gmail.com" // Замените на ваш email
		password = "yesz ypaf uxcd gmdz"        // Замените на ваш пароль приложения
		toEmail  = "cctrfsh@gmail.com"          // Email службы поддержки
	)

	m := gomail.NewMessage()

	// Отправитель
	m.SetHeader("From", username)
	// Получатель
	m.SetHeader("To", toEmail)
	// Тема письма
	m.SetHeader("Subject", subject)
	// Текст письма
	m.SetBody("text/plain", fmt.Sprintf("Message from: %s\n\n%s", fromEmail, body))

	// Добавление вложения, если файл был загружен
	if attachmentPath != "" {
		if _, err := os.Stat(attachmentPath); os.IsNotExist(err) {
			return fmt.Errorf("attachment not found: %s", attachmentPath)
		}
		m.Attach(attachmentPath)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, username, password)

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
