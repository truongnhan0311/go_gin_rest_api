package main

import (
	"github.com/chenyahui/gin-cache"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"warranty-api/config"
	"warranty-api/controllers"
	"warranty-api/middlewares"
	"warranty-api/services"
)

func main() {
	config.LoadEnv()
	router := setupRouter()
	router.Static("/api/v1/download", "./download")

	router.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})
	_ = router.Run(":" + os.Getenv("SERVER_PORT"))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.RateLimitMiddleware("20-S"))
	router.Use(middlewares.CORSMiddleware())
	router.GET("ping",
		cache.CacheByRequestURI(services.RedisStore(), 30*time.Second),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})

	v1 := router.Group("/api/v1")
	{
		warrantyRepo := controllers.WarrantyController()
		v1.GET("/search",
			cache.CacheByRequestURI(services.RedisStore(), 300*time.Second),
			warrantyRepo.SearchWarranty)
		v1.GET("/search/detail/:id",
			cache.CacheByRequestURI(services.RedisStore(), 300*time.Second),
			warrantyRepo.GetWarrantyByCode)
		v1.POST("/download/:id",
			cache.CacheByRequestURI(services.RedisStore(), 300*time.Second),
			warrantyRepo.DownloadWarranty)

		UserController := controllers.UserController()
		v1.POST("/login", UserController.Login)
		masterUser := v1.Group("/masteruser")
		{
			masterUser.POST("/signup", UserController.CreateUser)
		}

		user := v1.Group("/user")
		user.Use(middlewares.Authenticate())
		{
			user.POST("/signup", UserController.CreateUser)
			user.GET("/all", UserController.GetUsers)
			user.GET("/:id", UserController.GetUser)
			user.GET("", UserController.GetMe)
			user.PUT("", UserController.UpdateMe)
			user.PUT("/password", UserController.UpdatePassword)
			user.PUT("/:id", UserController.UpdateUser)
			user.DELETE("/:id", UserController.DeleteUser)
		}

		warranty := v1.Group("/warranty")
		warranty.Use(middlewares.Authenticate())
		{
			warranty.POST("", warrantyRepo.CreateWarranty)
			warranty.GET("/list", warrantyRepo.GetWarrantyList)
			warranty.GET("/download", warrantyRepo.ExportWarrantyList)
			warranty.GET("/:id", warrantyRepo.GetWarranty)
			warranty.PUT("/:id", warrantyRepo.UpdateWarranty)
			warranty.DELETE("/:agency/:daily/:id", warrantyRepo.DeleteWarranty)
		}

		agencyRepo := controllers.AgencyController()
		agency := v1.Group("/agency")
		agency.Use(middlewares.Authenticate())
		{
			agency.POST("", agencyRepo.Create)
			agency.GET("/list", agencyRepo.GetList)
			agency.GET("/:id", agencyRepo.Get)
			agency.PUT("/:id", agencyRepo.Update)
			agency.DELETE("/:id", agencyRepo.Delete)
		}

		dailyRepo := controllers.DailyController()
		daily := v1.Group("/daily")
		daily.Use(middlewares.Authenticate())
		{
			daily.POST("", dailyRepo.Create)
			daily.GET("/list", dailyRepo.GetList)
			daily.GET("/:id", dailyRepo.Get)
			daily.PUT("/:id", dailyRepo.Update)
			daily.DELETE("/:id", dailyRepo.Delete)
		}
	}

	return router
}
