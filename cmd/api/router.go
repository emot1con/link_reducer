package api

import (
	"go_link_reducer/services/url"
	"go_link_reducer/services/user"
	"go_link_reducer/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewRoutes(DB *gorm.DB) *gin.Engine {
	router := gin.Default()
	validator := validator.New()

	userRepo := user.NewUserRepositoryImpl(DB)
	DB.AutoMigrate(&types.User{})
	userHandler := user.NewHandler(userRepo, validator)
	userHandler.RegisterRoute(router)

	URLRepo := url.NewURLRepositoryImpl(DB)
	URLHandler := url.NewURLHandler(URLRepo, validator)
	URLHandler.RegisterRoute(router)

	return router
}
