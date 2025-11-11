package config

import (
	"back-go/services/email"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Dsn struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}

func (d Dsn) String() string {
	return "host=" + d.Host + " user=" + d.User + " password=" + d.Password + " dbname=" + d.Dbname + " port=" + d.Port + " " + " sslmode=disable"
}

type Config struct {
	Port              int          `yaml:"port"`
	Dsn               Dsn          `yaml:"dsn"`
	Smtp              email.Config `yaml:"smtp"`
	GeolocationAPIKey string       `yaml:"geolocation_api_key"`
	JWTSecretKey      string       `yaml:"jwt_secret_key"`
	EmailDev          bool         `yaml:"email_dev"`
	PricingDev        bool         `yaml:"pricing_dev"`
}

func ReadConfig(file *os.File) (configuration Config, err error) {
	err = yaml.NewDecoder(file).Decode(&configuration)
	if err != nil {
		return
	}
	if configuration.Port == 0 {
		err = errors.New("configuration port is missing")
		return
	}

	if configuration.Dsn.Host == "" ||
		configuration.Dsn.User == "" ||
		configuration.Dsn.Port == "" ||
		configuration.Dsn.Dbname == "" ||
		configuration.Dsn.Password == "" {
		err = errors.New("dsn configuration has missing values")
		return
	}

	if (configuration.Smtp.SMTPHost == "" ||
		configuration.Smtp.SMTPPort == 0 ||
		configuration.Smtp.Sender == "" ||
		configuration.Smtp.Username == "" ||
		configuration.Smtp.Password == "") &&
		configuration.EmailDev == false {
		err = errors.New("smtp configuration has missing values, set email_dev to true if an smtp configuration cannot be provided")
		return
	}

	if configuration.PricingDev == false && configuration.GeolocationAPIKey == "" {
		err = errors.New("geolocation_api_key is missing, set pricing_dev to true if an api key cannot be provided")
		return
	}
	if configuration.JWTSecretKey == "" {
		err = errors.New("jwt_secret_key is missing, generate a 256 bit key in hex format")
		return
	}
	if len(configuration.JWTSecretKey) != 64 {
		err = errors.New("jwt secret key length must be 256 bits long (64 chars hex encode)")
		return
	}
	return
}
