package models

import (
	"PlantApp/database"
	"time"
)

//	type User struct {
//		Id        uint   `gorm:"primaryKey"`
//		NickName  string `gorm:"size(256)"`
//		Email     string `gorm:"size(256)"`
//		Password  string `gorm:"size(256)"`
//		CreatedAt time.Time
//		UpdatedAt time.Time
//	}
type User struct {
	Id        uint   `gorm:"primaryKey"`
	NickName  string `gorm:"size(256)"`
	Email     string `gorm:"size(256)"`
	Password  string `gorm:"size(256)"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
type UserToken struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	UserId    uint
	Token     string `gorm:"size(256)"`
	User      User
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func init() {
	db := database.GetDB()
	db.AutoMigrate(&User{}, &UserToken{})
}
