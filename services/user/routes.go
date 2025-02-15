package user

import (
	"fmt"
	"go_link_reducer/contract"
	"go_link_reducer/services/auth"
	"go_link_reducer/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	repo      contract.UserRepository
	validator *validator.Validate
}

func NewHandler(repo contract.UserRepository, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		repo:      repo,
		validator: validator,
	}
}

func (u *UserHandler) RegisterRoute(route *gin.Engine) {
	route.POST("/auth/register", u.Register)
	route.POST("/auth/login", u.Login)
}

func (u *UserHandler) Register(c *gin.Context) {
	var payload types.RegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := u.validator.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request %v", err)})
		return
	}

	if _, err := u.repo.GetUserByEmail(payload.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	hashPassword, err := auth.GenerateHashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to hash password: %v", err)})
		return
	}
	payload.Password = hashPassword

	result, err := u.repo.Register(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to register user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

func (u *UserHandler) Login(c *gin.Context) {
	var payload types.LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := u.validator.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request %v", err)})
		return
	}

	user, err := u.repo.GetUserByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := auth.ComparePasswords(user.Password, []byte(payload.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	// accesToken, err := auth.GenerateJWTToken(user.ID, "user", time.Hour*1)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to generate token: %v", err)})
	// 	return
	// }

	// refreshToken, err := auth.GenerateJWTToken(user.ID, "user", time.Hour*24*30*3)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to generate token: %v", err)})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		// "accesToken":      accesToken,
		// "refreshToken":    refreshToken,
		"exp":             time.Now().Add(time.Hour * 1).String(),
		"expRefreshToken": time.Now().Add(time.Hour * 24 * 30 * 3).String(),
		"role":            "user",
	})
}
