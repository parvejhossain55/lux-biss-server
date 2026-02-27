package email

import (
	"fmt"
	"net/smtp"
)

type Sender interface {
	SendEmail(to []string, subject string, body string) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMTPSender struct {
	config *SMTPConfig
}

func NewSMTPSender(cfg *SMTPConfig) *SMTPSender {
	return &SMTPSender{config: cfg}
}

func (s *SMTPSender) SendEmail(to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.From, to, msg)
}
