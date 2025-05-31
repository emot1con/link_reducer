package url

import (
	"fmt"
	"go_link_reducer/contract"
	"go_link_reducer/types"
	"log"
	"math/rand"
	"net/http"
	"strings"
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
	route.POST("/new/url", u.CreateURL)
	route.GET("/urls", u.GetAll)
	route.GET("/:short_code", u.Redirect)
	route.POST("/cron/delete-expired", u.Delete)
}

func (u *URLHandler) CreateURL(c *gin.Context) {
	log.Printf("Received %s request to %s from %s", c.Request.Method, c.Request.URL.Path, c.Request.Host)

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

	// Set expiration date to 7 days from now if not provided
	if payload.ExpirationDate.IsZero() {
		payload.ExpirationDate = time.Now().Add(7 * 24 * time.Hour)
	}
	result, err := u.repo.Create(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Get the base URL based on environment
	baseURL := getBaseURL(c)
	shortURL := fmt.Sprintf("%s/%s", baseURL, result.ShortCode)

	c.JSON(http.StatusOK, gin.H{
		"id":              result.ID,
		"original_url":    result.OriginalURL,
		"short_url":       shortURL,
		"created_at":      result.CreatedAt,
		"hit_count":       result.HitCount,
		"expiration_date": result.ExpirationDate,
	})
}

func (u *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("short_code")

	URL, err := u.repo.GetOne(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	if strings.Contains(userAgent, "Googlebot") || strings.Contains(userAgent, "bingbot") {
		log.Println("Prefetch detected, not counting")
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusFound, URL.OriginalURL)

	if err := u.repo.Update(URL.ID, URL.HitCount+1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (u *URLHandler) GetAll(c *gin.Context) {

	URLsData, err := u.repo.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, URLsData)
}

func (u *URLHandler) Delete(c *gin.Context) {
	if err := u.repo.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Deleted expired URLs")
}

func generateShortURL(length int) string {
	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(shortCode)
}

// getBaseURL returns the appropriate base URL based on the environment
func getBaseURL(c *gin.Context) string {
	// Check if we're running on Railway (or any production environment)
	if c.Request.Host != "" && c.Request.Host != "localhost:8080" && c.Request.Host != "127.0.0.1:8080" {
		// Use the actual host from the request with HTTPS (Railway enforces HTTPS)
		if strings.Contains(c.Request.Host, "railway.app") || strings.Contains(c.Request.Host, ".up.railway.app") {
			return "https://" + c.Request.Host
		}
		// For other production environments, also use HTTPS
		return "https://" + c.Request.Host
	}

	// For localhost development
	return "http://" + c.Request.Host
}
