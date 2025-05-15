package services

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Dialer *gomail.Dialer
}

func NewMailer() *Mailer {
	mailerHost := os.Getenv("MAILER_HOST")
	mailerPort, _ := strconv.Atoi(os.Getenv("MAILER_PORT"))
	mailerUsername := os.Getenv("MAILER_USERNAME")
	mailerPass := os.Getenv("MAILER_PASSWORD")
	d := gomail.NewDialer(mailerHost, mailerPort, mailerUsername, mailerPass)
	return &Mailer{Dialer: d}
}

func (m *Mailer) Send(from, to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.Dialer.DialAndSend(msg)
}

