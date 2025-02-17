package api

import (
	"go_link_reducer/services/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewRoutes(DB *gorm.DB) *gin.Engine {
	router := gin.Default()
	validator := validator.New()

	URLRepo := url.NewURLRepositoryImpl(DB)
	URLHandler := url.NewURLHandler(URLRepo, validator)
	URLHandler.RegisterRoute(router)

	return router
}
