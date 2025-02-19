package url

import (
	"github.com/gin-gonic/gin"
	"go_link_reducer/types"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type URLRepositoryImpl struct {
	DB *gorm.DB
}

func NewURLRepositoryImpl(DB *gorm.DB) *URLRepositoryImpl {
	return &URLRepositoryImpl{
		DB: DB,
	}
}

func (u *URLRepositoryImpl) Create(URLPayload types.CreateURLPayload) (types.URL, error) {
	URL := types.URL{
		OriginalURL:    URLPayload.OriginalURL,
		ShortCode:      URLPayload.ShortCode,
		HitCount:       0,
		ExpirationDate: URLPayload.ExpirationDate,
	}

	if err := u.DB.Create(&URL).Error; err != nil {
		return types.URL{}, err
	}

	return URL, nil
}

func (u *URLRepositoryImpl) GetOne(link string) (types.URL, error) {
	var URL types.URL
	if err := u.DB.Where("short_code = ?", link).First(&URL).Error; err != nil {
		return types.URL{}, err
	}

	return URL, nil
}

func (u *URLRepositoryImpl) GetAll(c *gin.Context) (map[string]any, error) {
	var URLs []types.URL

	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	u.DB.Model(URLs).Count(&total)

	offset := (page - 1) * 10

	totalPage := int(math.Ceil(float64(total) / float64(10)))

	if err := u.DB.Limit(10).Offset(offset).Find(&URLs).Error; err != nil {
		return nil, err
	}

	return map[string]any{
		"total":      total,
		"page":       page,
		"total_page": totalPage,
		"data":       URLs,
	}, nil
}

func (u *URLRepositoryImpl) Update(ID uint, count int) error {
	if err := u.DB.Model(&types.URL{}).Where("id = ?", ID).Update("hit_count", count).Error; err != nil {
		return err
	}
	return nil
}

func (u *URLRepositoryImpl) Delete() error {
	if err := u.DB.Where("expiration_date < ?", time.Now()).Delete(&types.URL{}).Error; err != nil {
		return err
	}

	return nil
}
