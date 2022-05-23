package email

import (
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

type Email struct {
	userName string
}

func NewEmailClient(username string) *Email {
	return &Email{
		userName: username,
	}
}

func (e *Email) SendEmail(email string, msgBody []byte) error {
	from := e.userName
	password := viper.GetString("app.emailPassword")

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Go final")
	msg.SetBody("text/html", "Your email verificatin link:\n"+string(msgBody))

	n := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
