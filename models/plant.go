package models

import (
	"PlantApp/database"
	"time"
)

type Plant struct {
	Id          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size(256)"`
	Image       string `gorm:"size(256)"`
	Description string `gorm:"size(256)"`
	Uses        string `gorm:"size(256)"`
	Soil        string `gorm:"size(256)"`
	Climate     string `gorm:"size(256)"`
	Health      string `gorm:"size(256)"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type PlantUser struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint
	PlantId   uint
	Name      string `gorm:"size(256)"`
	Image     string `gorm:"size(512)"`
	User      User   `gorm:"foreignKey:UserId"`
	Plant     Plant  `gorm:"foreignKey:PlantId"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Rose struct {
	Id          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size(256)"`
	Image       string `gorm:"size(256)"`
	Description string `gorm:"size(256)"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func init() {
	db := database.GetDB()
	db.AutoMigrate(&Plant{}, &Rose{}, &PlantUser{})
}
