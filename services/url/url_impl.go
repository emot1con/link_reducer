package url

import (
	"go_link_reducer/types"

	"gorm.io/gorm"
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
		OriginalURL: URLPayload.OriginalURL,
		ShortCode:   URLPayload.ShortCode,
		HitCount:    0,
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

func (u *URLRepositoryImpl) GetAll() ([]types.URL, error) {
	var URLs []types.URL
	if err := u.DB.Find(&URLs).Error; err != nil {
		return nil, err
	}

	return URLs, nil
}

func (u *URLRepositoryImpl) Update(ID uint, count int) error {
	if err := u.DB.Model(&types.URL{}).Where("id = ?", ID).Update("hit_count", count).Error; err != nil {
		return err
	}
	return nil
}

func (u *URLRepositoryImpl) Delete(URL string) error {
	if err := u.DB.Where("short_code = ?", URL).Delete(&types.URL{}).Error; err != nil {
		return err
	}

	return nil
}
