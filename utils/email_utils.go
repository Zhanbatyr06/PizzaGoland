package utils

import (
	"log"

	"github.com/go-gomail/gomail"
)

// SendEmail отправляет email
func SendEmail(to, subject, body string, attachmentPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "zhanbatyrmolkryt@gmail.com") // Ваша почта
	m.SetHeader("To", to)                             // Почта получателя
	m.SetHeader("Subject", subject)                   // Тема
	m.SetBody("text/plain", body)                     // Тело письма

	if attachmentPath != "" {
		m.Attach(attachmentPath)
	}

	// Настройка SMTP
	d := gomail.NewDialer("smtp.gmail.com", 587, "zhanbatyrmolkryt@gmail.com", "yesz ypaf uxcd gmdz")

	// Отправка
	log.Printf("Sending email to: %s with subject: %s", to, subject)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
