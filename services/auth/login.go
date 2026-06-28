package auth

import (
	"Trading-Risk_Management/infrastructure/jwt"
	"Trading-Risk_Management/models"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (s *AuthServiceImpl) Login(username, password string) (*AuthTokens, error) {
	existing, err := s.userRepo.GetByUsername(username)
	if existing == nil || err != nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	accessToken, err := jwt.Token(existing, s.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &models.RefreshToken{
		UserID:    existing.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	s.refreshTokenRepo.Create(rt)

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

func (s *AuthServiceImpl) RefreshToken(refreshToken string) (*AuthTokens, error) {
	rt, err := s.refreshTokenRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		s.refreshTokenRepo.Revoke(refreshToken)
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.GetByID(rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	s.refreshTokenRepo.Revoke(refreshToken)

	accessToken, err := jwt.Token(user, s.cfg)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRT := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	s.refreshTokenRepo.Create(newRT)

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    900,
	}, nil
}

func (s *AuthServiceImpl) Logout(refreshToken string) error {
	return s.refreshTokenRepo.Revoke(refreshToken)
}
