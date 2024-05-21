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
	router.GET("/deneme", controllers.RequireAuth(), controllers.GetImage)
	mobil := router.Group("/mobil")
	{
		mobil.POST("/user/create", controllers.UserRegister)
		mobil.POST("/user/login", controllers.Login)
		mobil.Use(controllers.RequireAuth())
		{
			mobil.GET("/user/logout", controllers.Logout)
			mobil.GET("/user", controllers.GetUserById)
			mobil.GET("/myplants", controllers.MyPlants)
			mobil.GET("/rose", controllers.Rose)
			mobil.GET("/leaf", controllers.Leaf)
			mobil.GET("/plants", controllers.GetPlant)
		}
	}

}
