package auth

import (
	"Trading-Risk_Management/config"
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	cfg              *config.Config
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, cfg *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo, cfg: cfg}
}

func (s *AuthServiceImpl) Register(user models.User) error {
	existing, _ := s.userRepo.GetByEmail(user.Email)
	if existing != nil {
		return errors.New("duplicate email")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.userRepo.Create(&user)
}

func (s *AuthServiceImpl) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("wrong password")
	}

	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return s.userRepo.Update(user)
}
