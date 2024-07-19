package mail

import (
	"emailn/internal/domain/campaign"
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendMail(campaign *campaign.Campaign) error {

	setupMail := gomail.NewDialer("smtp.gmail.com", 587, "test@gmail.com", "123456")

	var emails []string
	for _, conctact := range campaign.Contacts {
		emails = append(emails, conctact.Email)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "test@gmail.com")
	message.SetHeader("TO", emails...)
	message.SetAddressHeader("Cc", "suporte@dev.com.br", "Dan")
	message.SetHeader("Subject", campaign.Name)
	message.SetBody("text/html", campaign.Content)

	if err := setupMail.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
