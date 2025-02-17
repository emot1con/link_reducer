package types

import (
	"time"
)

type URL struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OriginalURL    string    `gorm:"not null" json:"original_url"`
	ShortCode      string    `gorm:"unique;not null" json:"short_code"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	HitCount       int       `gorm:"default:0" json:"hit_count"`
	ExpirationDate time.Time `gorm:"" json:"expiration_date"`
}

type CreateURLPayload struct {
	OriginalURL    string    `json:"original_url" validate:"required,url"`
	ShortCode      string    `json:"short_code,omitempty"`
	ExpirationDate time.Time `json:"expiration_date,omitempty"`
}
