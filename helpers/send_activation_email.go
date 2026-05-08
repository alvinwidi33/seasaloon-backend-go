package helpers

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendActivationEmail(toEmail, token string) error {
	from := os.Getenv("EMAIL_HOST_USER")
	password := os.Getenv("EMAIL_HOST_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	activationLink := fmt.Sprintf(
		"http://localhost:8081/api/activate?token=%s",
		token,
	)

	htmlBytes, err := os.ReadFile("helpers/email/activation.html")
	if err != nil {
		return err
	}

	templateHTML := string(htmlBytes)

	html := strings.ReplaceAll(templateHTML, "{{NAME}}", toEmail)
	html = strings.ReplaceAll(html, "{{LINK}}", activationLink)

	subject := "Subject: Verifikasi Akun\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	message := []byte(subject + mime + html)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toEmail},
		message,
	)
}