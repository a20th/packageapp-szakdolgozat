package email

import (
	"io"

	"gopkg.in/gomail.v2"
)

type consoleImpl struct {
	Writer io.Writer
}

func (c consoleImpl) SendEmail(message gomail.Message) error {
	_, err := message.WriteTo(c.Writer)
	if err != nil {
		return err
	}
	return nil
}

func CreateConsoleService(writer io.Writer) Service {
	return consoleImpl{writer}
}
