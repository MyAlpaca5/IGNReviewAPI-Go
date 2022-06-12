package models

import (
	"time"
)

type Review struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" binding:"required"`
	UpdatedAt   time.Time `json:"updated_at" binding:"required,gtefield=CreatedAt"`
	ReviewURL   string    `json:"review_url" binding:"required,ignurl"`
	ReviewScore float32   `json:"review_score" binding:"required,min=0,max=10"`
	MediaType   string    `json:"media_type,omitempty"`
	GenreList   []string  `json:"genre_list,omitempty"`
	CreatorList []string  `json:"creator_list,omitempty"`
}
