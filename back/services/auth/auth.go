package auth

import (
	"back-go/services/account"
	"back-go/services/models"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	stdjwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	Login(email string, password string) (string, string, error)
	Logout(id string) error
	Refresh(email string, id string) (string, string, error)
}

type service struct {
	JWTKey      []byte
	AccountRepo account.Repository
	AuthRepo    Repository
}

func createNewClaims(email string) (accessClaims stdjwt.RegisteredClaims, refreshClaims stdjwt.RegisteredClaims) {
	accessExpiration := stdjwt.NewNumericDate(time.Now().Add(time.Hour))
	refreshExpiration := stdjwt.NewNumericDate(time.Now().Add(time.Hour * 168))
	accessClaims = stdjwt.RegisteredClaims{
		Subject:   email,
		Audience:  stdjwt.ClaimStrings{"access", "user"},
		ExpiresAt: accessExpiration,
		IssuedAt:  stdjwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}

	refreshClaims = stdjwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: refreshExpiration,
		Audience:  stdjwt.ClaimStrings{"refresh", "user"},
		IssuedAt:  stdjwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}

	return
}

func (s service) Login(email string, password string) (accessToken string, refreshToken string, err error) {
	acc, err := s.AccountRepo.FindByEmail(email)
	if errors.Is(err, errors.New("no records found")) {
		err = ErrInvalidCredentials
		return
	}
	if err != nil {
		return
	}

	if acc.VerifiedAt == nil {
		err = account.ErrNotVerified
		return
	}

	match, err := argon2id.ComparePasswordAndHash(password, acc.Password)
	if err != nil {
		return
	}

	if !match {
		err = ErrInvalidCredentials
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

func CreateAuthService(JWTKey []byte, AccountRepo account.Repository, AuthRepo Repository) (s Service) {
	return &service{
		JWTKey:      JWTKey,
		AccountRepo: AccountRepo,
		AuthRepo:    AuthRepo,
	}
}

type Repository interface {
	Store(tokens *models.TokenRecord) error
	Find(id string) (*models.TokenRecord, error)
	Delete(id string) error
}
