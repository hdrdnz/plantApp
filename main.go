package main

import (
	"PlantApp/routers"

	"github.com/gin-gonic/gin"
)

// @title          Leaflove Mobil API
// @version         1.0
// @description     This is plantapp

// @host      https://leaflove.com.tr
func main() {
	router := gin.Default()
	router.Static("./images", "images")
	router.Static("./build", "/")
	routers.Load(router)
	router.Run()

}
