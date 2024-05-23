package routers

import (
	controllers "PlantApp/controllers/user"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Load(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL("/getdocs")))
	router.GET("/getdocs", controllers.GetDocs)
	router.POST("/deneme", controllers.RequireAuth(), controllers.GetImage)
	router.POST("/getimage", controllers.Images)
	router.POST("/postplant", controllers.PostPlant)
	router.POST("/post-fruits", controllers.GeneralPlant)
	//router.GET("/favorites", controllers.GetFavorites)

	mobil := router.Group("/mobil")
	{
		mobil.POST("/user/create", controllers.UserRegister)
		mobil.POST("/user/login", controllers.Login)
		mobil.Use(controllers.RequireAuth())
		{
			mobil.GET("/user/logout", controllers.Logout)
			mobil.GET("/user", controllers.GetUserById)
			//bitkilerim kısmı
			mobil.GET("/plants", controllers.MyPlants)
			//genel bilgi için
			mobil.GET("/fruits", controllers.GetFruits)
			//favoriler kısmı
			mobil.GET("/favorites", controllers.GetFavorites)
			mobil.GET("/rose", controllers.Rose)
			mobil.GET("/leaf", controllers.Leaf)
			mobil.GET("/plant", controllers.GetPlant)
		}
	}

}
