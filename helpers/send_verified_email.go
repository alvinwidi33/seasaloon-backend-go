package helpers
import (
	"net/smtp"
	"os"
	"strings"
)
func SendVerifiedEmail(toEmail string) error {
	from := os.Getenv("EMAIL_HOST_USER")
	password := os.Getenv("EMAIL_HOST_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	htmlBytes, err := os.ReadFile("helpers/email/verified.html")
	if err != nil {
		return err
	}

	templateHTML := string(htmlBytes)

	frontendURL := os.Getenv("FRONTEND_URL")
	loginURL := frontendURL + "/login"

	html := strings.ReplaceAll(templateHTML, "{{LOGIN_URL}}", loginURL)

	subject := "Subject: Akun Berhasil Diverifikasi\r\n"
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