package models

import "gorm.io/gorm"

type Menfest struct {
	gorm.Model
	Pesan  string
	UserID string
	Users  []User `gorm:"many2many:user_menfest"`
	// UserIDs []uint `gorm:"-"`
}
