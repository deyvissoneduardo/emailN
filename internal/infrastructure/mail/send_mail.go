package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendMail() error {

	setupMail := gomail.NewDialer("smtp.gmail.com", 587, "test@gmail.com", "123456")

	message := gomail.NewMessage()
	message.SetHeader("From", "test@gmail.com")
	message.SetHeader("TO", "test@gmail.com")
	message.SetAddressHeader("Cc", "test@gmail.com", "Dan")
	message.SetHeader("Subject", "Hello!")
	message.SetBody("text/html", "Hello <p1>Deyvisson</p1>")

	if err := setupMail.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
