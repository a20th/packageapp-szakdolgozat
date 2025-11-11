package admin

import (
	"back-go/services/auth"
	"back-go/services/models"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	stdjwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func createNewClaims(email string) (accessClaims stdjwt.RegisteredClaims, refreshClaims stdjwt.RegisteredClaims) {
	accessExpiration := stdjwt.NewNumericDate(time.Now().Add(time.Hour))
	refreshExpiration := stdjwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	accessClaims = stdjwt.RegisteredClaims{
		Subject:   email,
		Audience:  stdjwt.ClaimStrings{"access", "admin"},
		ExpiresAt: accessExpiration,
		IssuedAt:  stdjwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}

	refreshClaims = stdjwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: refreshExpiration,
		Audience:  stdjwt.ClaimStrings{"refresh", "admin"},
		IssuedAt:  stdjwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}

	return
}

type Service interface {
	Login(string, string) (string, string, error)
	Logout(string) error
	Refresh(string, string) (string, string, error)
	Create(string, string) error
	Delete(string) error
	GetUsers() ([]string, error)
}

func CreateAdminService(jwtKey []byte, adminRepo Repository, authRepo auth.Repository) Service {
	return &service{JWTKey: jwtKey, Repo: adminRepo, AuthRepo: authRepo}
}

type service struct {
	JWTKey   []byte
	Repo     Repository
	AuthRepo auth.Repository
}

func (s service) GetUsers() (st []string, err error) {
	users, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	st = make([]string, len(users))
	for i, user := range users {
		st[i] = user.Username
	}
	return
}

func (s service) Delete(username string) error {
	err := s.Repo.Delete(username)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Create(username string, password string) error {

	pwd, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	admin := models.Admin{Username: username, Password: pwd}
	err = s.Repo.Store(&admin)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Login(id string, password string) (accessToken string, refreshToken string, err error) {
	acc, err := s.Repo.Find(id)
	if errors.Is(err, errors.New("no records found")) {
		err = auth.ErrInvalidCredentials
		return
	}
	if err != nil {
		return
	}

	match, err := argon2id.ComparePasswordAndHash(password, acc.Password)
	if err != nil {
		return
	}

	if !match {
		err = auth.ErrInvalidCredentials
		return
	}

	accessClaims, refreshClaims := createNewClaims(id)

	tokenRecord := models.TokenRecord{
		AccessID:  accessClaims.ID,
		RefreshID: refreshClaims.ID,
	}

	access := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, accessClaims)
	refresh := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, refreshClaims)
	accessToken, err = access.SignedString(s.JWTKey)
	if err != nil {
		return
	}
	refreshToken, err = refresh.SignedString(s.JWTKey)
	if err != nil {
		return
	}

	err = s.AuthRepo.Store(&tokenRecord)
	if err != nil {
		return
	}
	return
}

func (s service) Logout(id string) (err error) {
	err = s.AuthRepo.Delete(id)
	return
}

func (s service) Refresh(id string, email string) (accessToken string, refreshToken string, err error) {
	token, err := s.AuthRepo.Find(id)
	if token == nil || token.RefreshID != id {
		return
	}
	err = s.AuthRepo.Delete(id)
	if err != nil {
		return
	}

	accessClaims, refreshClaims := createNewClaims(email)

	tokenRecord := models.TokenRecord{
		AccessID:  accessClaims.ID,
		RefreshID: refreshClaims.ID,
	}

	access := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, accessClaims)
	refresh := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, refreshClaims)
	accessToken, err = access.SignedString(s.JWTKey)
	if err != nil {
		return
	}
	refreshToken, err = refresh.SignedString(s.JWTKey)
	if err != nil {
		return
	}

	err = s.AuthRepo.Store(&tokenRecord)
	if err != nil {
		return
	}
	return
}

type Repository interface {
	Store(*models.Admin) error
	Find(string) (*models.Admin, error)
	GetAll() ([]models.Admin, error)
	Delete(string) error
}
