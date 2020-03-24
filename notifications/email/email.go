package email

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmail(to string, msg string) {
	auth := smtp.PlainAuth("", os.Getenv("FROM_EMAIL"), os.Getenv("PASSWORD"), os.Getenv("SMTP_SERVER"))

	err := smtp.SendMail(os.Getenv("SMTP_SERVER")+":25", auth, os.Getenv("FROM_EMAIL"), []string{to}, []byte(msg))

	if err != nil {
		log.Fatalf("Failed to send email: %s", err)
	}
}
