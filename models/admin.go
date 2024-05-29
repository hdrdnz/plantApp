package models

import (
	"PlantApp/database"
	"time"
)

type Admin struct {
	Id        uint   `gorm:"primaryKey"`
	NickName  string `gorm:"size(256)"`
	Email     string `gorm:"size(256)"`
	Password  string `gorm:"size(256)"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type AdminToken struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	AdminId   uint
	Token     string `gorm:"size(256)"`
	Admin     Admin
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func init() {
	db := database.GetDB()
	db.AutoMigrate(&Admin{}, &AdminToken{})
}
