package email

import (
	"crypto/tls"

	gomail "gopkg.in/gomail.v2"
)

type Service interface {
	SendEmail(message gomail.Message) error
}

type service struct {
	config Config
}

type Config struct {
	SMTPHost string `json:"host"`
	SMTPPort int    `json:"port"`
	Sender   string
	Username string
	Password string
}

func (s service) SendEmail(message gomail.Message) error {
	message.SetHeader("From", s.config.Sender)
	dialer := dialFromConfig(s.config)

	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
	}

	sender, err := dialer.Dial()
	if err != nil {
		return err
	}

	err = gomail.Send(sender, &message)
	if err != nil {
		return err
	}

	return nil
}

func dialFromConfig(s Config) *gomail.Dialer {
	return gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.Username, s.Password)
}

func CreateEmailService(config Config) Service {
	return &service{config: config}
}

func TestEmailConfig(config Config) (err error) {
	dialer := dialFromConfig(config)
	c, err := dialer.Dial()
	c.Close()
	return
}
