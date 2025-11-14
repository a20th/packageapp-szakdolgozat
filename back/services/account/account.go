package account

import (
	"back-go/services/email"
	"back-go/services/models"
	"errors"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-passwd/validator"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var ErrAlreadyVerified = errors.New("already verified")
var ErrNotVerified = errors.New("account not verified")
var ErrInvalidPhoneNumber = errors.New("invalid phone number")
var ErrInvalidVerificationCode = errors.New("invalid verification code")
var ErrAccountNotExist = errors.New("account does not exist")
var ErrInvalidEmail = errors.New("invalid email")
var ErrDuplicateEmail = errors.New("duplicate email")

type Service interface {
	Register(email string, password string, name string, phoneNumber string, preferredLang string) error
	Get(email string) (*models.Account, error)
	Verify(email string, code string) error
	Update(account models.Account) error
	Delete(email string) error
}

type service struct {
	Repo             Repository
	EmailService     email.Service
	FrontendLocation string
}

func (s service) Register(email string, password string, name string, phoneNumber string, preferredLang string) error {

	if preferredLang != "hu" && preferredLang != "en" {
		return ErrInvalidPhoneNumber
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	num, err := phonenumbers.Parse(phoneNumber, "HU")
	if err != nil {
		return ErrInvalidPhoneNumber
	}
	if *num.CountryCode != 36 {
		return ErrInvalidPhoneNumber
	}
	phoneNumber = phonenumbers.Format(num, phonenumbers.E164)

	passwordValidator := validator.New(
		validator.MaxLength(32, nil),
		validator.MinLength(8, nil),
		validator.ContainsAtLeast("0123456789", 1, nil),
		validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 1, nil))
	err = passwordValidator.Validate(password)
	if err != nil {
		return err
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	acc := &models.Account{
		Password:          hash,
		Email:             email,
		Name:              name,
		PhoneNumber:       phoneNumber,
		PreferredLanguage: preferredLang,
	}

	verificationCode := uuid.New().String()
	acc.VerificationID = &verificationCode

	var message gomail.Message
	{
		message = *gomail.NewMessage()
		message.SetHeader("To", email)

		site := "<h2>Kedves %name</h2><p>A felhasználó fiókod aktiválásához kattints az alábbi linkre:</p><a href=\"%link\">%link</a>" +
			"<p>To activate your account click on the link above:</p>"
		message.SetHeader("Subject", "Regisztráció megerősítése / Verify registration")

		//TODO
		link := "http://localhost:5173/app?verify=" + verificationCode + "&email=" + url.QueryEscape(email)

		site = strings.Replace(site, "%name", acc.Name, -1)
		site = strings.Replace(site, "%link", link, -1)

		message.SetBody("text/html", site)
	}
	err = s.EmailService.SendEmail(message)
	if err != nil {
		return err
	}

	err = s.Repo.Store(acc)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (s service) Verify(email string, code string) error {
	account, err := s.Repo.FindByEmail(email)
	if err != nil {
		return err
	}

	if *account.VerificationID != code {
		return ErrInvalidVerificationCode
	}
	t := time.Now()
	account.VerifiedAt = &t
	err = s.Repo.Store(account)
	if err != nil {
		return err
	}

	return nil
}

func (s service) Get(email string) (*models.Account, error) {
	var account *models.Account
	account, err := s.Repo.FindByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotExist
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s service) Update(account models.Account) error {
	_, err := s.Repo.FindByEmail(account.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotExist
	}

	err = s.Repo.Store(&account)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Delete(id string) error {
	_, err := s.Repo.Find(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotExist
	}
	return s.Repo.Delete(id)
}

func CreateAccountService(repo Repository, emailService email.Service) Service {
	return &service{
		Repo:         repo,
		EmailService: emailService,
	}
}

type Repository interface {
	Store(account *models.Account) error
	Delete(id string) error
	Find(id string) (*models.Account, error)
	FindByEmail(email string) (*models.Account, error)
}
