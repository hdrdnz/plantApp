package routers

import (
	admin "PlantApp/controllers/admin"
	user "PlantApp/controllers/user"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Load(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("/getdocs")))
	router.GET("/getdocs", user.GetDocs)
	adminr := router.Group("/admin")
	{
		adminr.POST("/register", admin.AdminRegister)
		adminr.POST("/login", admin.AdminLogin)
		adminr.Use(admin.RequireAuth())
		{
			adminr.GET("/plants", admin.PlantGet)
			adminr.POST("/add-plant", admin.PlantPost)

			adminr.POST("/plant/:plantid/update", admin.UpdatePlant)
			adminr.POST("/plant/:plantid/delete", admin.DeletePlant)
		}

	}
	router.GET("/web/plants", user.WebPlants)

	mobil := router.Group("/mobil")
	{
		mobil.POST("/user/create", user.UserRegister)
		mobil.POST("/user/login", user.Login)

		mobil.Use(user.RequireAuth())
		{
			mobil.GET("/user/logout", user.Logout)
			mobil.GET("/user", user.GetUserById)
			//bitkilerim kısmı
			mobil.GET("/plants", user.MyPlants)
			mobil.POST("/plant-upload", user.GetInformation)
			mobil.GET("/plant-detail", user.GetDetail)
			//genel bilgi için
			mobil.GET("/general-plants", user.GetFruits)
			//favoriler kısmı
			mobil.POST("/add-favorite", user.AddFavorite)
			mobil.POST("/delete-favorite", user.DelFavor)
			mobil.GET("/favorites", user.GetFavorites)
			mobil.GET("/rose", user.Rose)
			mobil.GET("/plant", user.GetPlant)
		}
	}

}
