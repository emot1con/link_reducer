package api

import (
	"go_link_reducer/services/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewRoutes(DB *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://link-reducer.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	router.Use(gin.Logger())

	validator := validator.New()

	URLRepo := url.NewURLRepositoryImpl(DB)
	URLHandler := url.NewURLHandler(URLRepo, validator)
	URLHandler.RegisterRoute(router)

	return router
}
