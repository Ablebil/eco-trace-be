package email

import (
	"github.com/Ablebil/eco-sample/config"
	"gopkg.in/gomail.v2"
)

type EmailItf interface {
	SendOTPEmail(to, otp string) error
}

type Email struct {
	sender   string
	password string
}

func NewEmail(cfg *config.Config) EmailItf {
	return &Email{
		sender:   cfg.EmailUser,
		password: cfg.EmailPassword,
	}
}

func (e *Email) SendOTPEmail(to, otp string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", e.sender)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Your OTP Code")
	mail.SetBody("text/plain", "Your OTP code is: "+otp)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, e.sender, e.password)
	return dialer.DialAndSend(mail)
}
