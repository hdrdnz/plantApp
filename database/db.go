package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	// gotoErr := godotenv.Load(".env")
	// if gotoErr != nil {
	// 	log.Fatal("err:", gotoErr)
	// 	return nil
	// }
	// // dbUser := os.Getenv("mysqluser")
	// dbPass := os.Getenv("mysqlpass")
	// dbHost := os.Getenv("mysqlhost")
	// dbName := os.Getenv("mysqldb")

	var DBError error
	dsn := "leafloveusr" + ":" + "F0!DqcRf5X0BE*(3" + "@tcp(" + "93.115.79.72" + ":3306)/" + "leaflove" + "?parseTime=true&charset=utf8&loc=UTC"
	DB, DBError = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if DBError != nil {
		log.Fatal("err:", DBError)
		return nil
	}
	return DB
}
func AutoMigrate(model *struct{}) {
	DB.AutoMigrate(model)
}
