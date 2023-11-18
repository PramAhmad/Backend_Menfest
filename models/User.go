package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid     string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	Username string
	Token    string
	Menfest  []Menfest `gorm:"many2many:user_menfest"`
}
