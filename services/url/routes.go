package url

import (
	"fmt"
	"go_link_reducer/contract"
	"go_link_reducer/types"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type URLHandler struct {
	repo      contract.URLRepository
	validator *validator.Validate
}

func NewURLHandler(repo contract.URLRepository, validator *validator.Validate) *URLHandler {
	return &URLHandler{
		repo:      repo,
		validator: validator,
	}
}

func (u *URLHandler) RegisterRoute(route *gin.Engine) {
	route.POST("/urls", u.CreateURL)
	route.GET("urls", u.GetAll)
	route.GET("/:short_code", u.Redirect)
	route.DELETE("/urls/:short_code", u.Delete)
}

func (u *URLHandler) CreateURL(c *gin.Context) {
	var payload types.CreateURLPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := u.validator.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if payload.ShortCode != "" {
		if _, err := u.repo.GetOne(payload.ShortCode); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "short code already exists"})
			return
		}
	} else {
		payload.ShortCode = generateShortURL(6)
		for _, err := u.repo.GetOne(payload.ShortCode); err == nil; {
			payload.ShortCode = generateShortURL(6)
		}
	}

	result, err := u.repo.Create(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", result.ShortCode)

	c.JSON(http.StatusOK, gin.H{
		"id":           result.ID,
		"original_url": result.OriginalURL,
		"short_url":    shortURL,
		"created_at":   result.CreatedAt,
		"hit_count":    result.HitCount,
	})
}

func (u *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("short_code")

	URL, err := u.repo.GetOne(shortCode)
	log.Println(URL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, URL.OriginalURL)
}

func (u *URLHandler) GetAll(c *gin.Context) {
	URLsData, err := u.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": URLsData})
}

func (u *URLHandler) Delete(c *gin.Context) {
	shortCode := c.Param("short_code")
	if err := u.repo.Delete(shortCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted"})
}

func generateShortURL(length int) string {
	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(shortCode)
}
