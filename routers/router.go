package routers

import (
	"mygram/controllers"
	"mygram/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer() *gin.Engine {
	router := gin.Default()
	usersRoute := router.Group("/users")
	{

		usersRoute.POST("/register", controllers.UserRegister)
		usersRoute.POST("/login", controllers.UserLogin)
		usersRoute.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", gin.H{})
		})
		usersRoute.Use(middlewares.Authentication())
		usersRoute.PUT("/", controllers.UserUpdate)
		usersRoute.DELETE("/", controllers.UserDelete)
	}
	photosRoute := router.Group("/photos")
	{
		photosRoute.GET("/", controllers.PhotoList)
		photosRoute.GET("/:id", controllers.PhotoListByID)
		photosRoute.Use(middlewares.Authentication())
		photosRoute.POST("/", controllers.PhotoUpload)
		photosRoute.PUT("/:id", controllers.PhotoUpdate)
		photosRoute.DELETE("/:id", controllers.PhotoDelete)

	}
	commentRoute := router.Group("/comments")
	{
		commentRoute.Use(middlewares.Authentication())
		commentRoute.POST("/", controllers.CommentPost)
	}

	router.Static("/public", "./templates")
	router.GET("/", func(c *gin.Context) {
		c.File("./templates/test.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	router.Use(Cors())
	return router
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
