package models

import (
	"time"
)

// Review is a struct used in POST method and also used for passing review to different method.
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
	Version     int64     `json:"-"`
}

// ReviewForUpdate is used in PATCH method.
// Using pointer here to differentiate from user input zero/default value with not inputing any value for a field.
type ReviewForUpdate struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" binding:"required"`
	ReviewScore *float32  `json:"review_score" binding:"omitempty,min=0,max=10"`
	MediaType   *string   `json:"media_type,omitempty"`
	GenreList   []string  `json:"genre_list,omitempty"`
	CreatorList []string  `json:"creator_list,omitempty"`
}
