package mailer

import (
	"fmt"
	"net/smtp"

	"github.com/slipe-fun/skid-backend/internal/config"
)

func SendMail(subject string, content string, recipient string) error {
	cfg := config.LoadConfig("configs/config.yaml")

	to := []string{recipient}

	headerMap := make(map[string]string)
	headerMap["From"] = cfg.Email.Email
	headerMap["To"] = to[0]
	headerMap["Subject"] = subject
	headerMap["MIME-Version"] = "1.0"
	headerMap["Content-Type"] = "text/plain; charset=\"UTF-8\""

	messageStr := ""
	for k, v := range headerMap {
		messageStr += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	messageStr += "\r\n" + content
	message := []byte(messageStr)

	auth := smtp.PlainAuth("", cfg.Email.SmtpLogin, cfg.Email.SmtpPassword, cfg.Email.SmtpHost)

	err := smtp.SendMail(cfg.Email.SmtpHost+":"+cfg.Email.SmtpPort, auth, cfg.Email.Email, to, message)

	if err != nil {
		return err
	}

	return nil
}
