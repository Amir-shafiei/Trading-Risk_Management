package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"email" gorm:"unique"`
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"-"`
}
