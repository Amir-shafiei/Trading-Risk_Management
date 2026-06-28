package repository

import (
	"Trading-Risk_Management/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}
func (repo *UserRepo) Create(user *models.User) error {
	return repo.db.Create(user).Error
}
func (repo *UserRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := repo.db.Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}
func (repo *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}
func (repo *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := repo.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}

func (repo *UserRepo) Update(user *models.User) error {
	return repo.db.Save(user).Error
}
