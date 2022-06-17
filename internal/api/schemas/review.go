package schemas

import (
	"errors"
	"regexp"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/db/models"
)

// ReviewIn struct includes data that is received from a request
// Used pointer here to differentiate from zero value and no input given
type ReviewIn struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	ReviewScore *float32 `json:"review_score"`
	ReviewURL   *string  `json:"review_url"`
	MediaType   *string  `json:"media_type,omitempty"`
	GenreList   []string `json:"genre_list,omitempty"`
	CreatorList []string `json:"creator_list,omitempty"`
}

// FillReview accepts a Reviwe model and update it with data in a ReviewIn schema
func (r ReviewIn) UpdateReview(original models.Review) models.Review {
	var updated = original

	if r.Name != nil {
		updated.Name = *r.Name
	}
	if r.Description != nil {
		updated.Description = *r.Description
	}
	if r.ReviewURL != nil {
		updated.ReviewURL = *r.ReviewURL
	}
	if r.ReviewScore != nil {
		updated.ReviewScore = *r.ReviewScore
	}
	if r.MediaType != nil {
		updated.MediaType = *r.MediaType
	}
	if r.GenreList != nil {
		updated.GenreList = r.GenreList
	}
	if r.CreatorList != nil {
		updated.CreatorList = r.CreatorList
	}

	return updated
}

func (r *ReviewIn) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return errors.New("field 'name' cannot be empty string")
	}

	if r.ReviewScore != nil && (*r.ReviewScore < 0 || *r.ReviewScore > 10) {
		return errors.New("field 'review_score' must be a float and its values has been with [0, 10]")
	}

	if r.ReviewURL != nil {
		matched, err := regexp.MatchString("^https://www.ign.com/articles/[12][0-9]{3}/(0[1-9]|1[0-2])/(0[1-9]|[1-2][0-9]|3[01])/.+$", *r.ReviewURL)
		if !matched || err != nil {
			return errors.New("field 'review_url' must be a valid link to a IGN article")
		}
	}

	return nil
}

// ReviewOut struct includes data that is sent back via a response
type ReviewOut struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ReviewURL   string   `json:"review_url"`
	ReviewScore float32  `json:"review_score"`
	MediaType   string   `json:"media_type"`
	GenreList   []string `json:"genre_list"`
	CreatorList []string `json:"creator_list"`
}

func ToReviewOut(r models.Review) ReviewOut {
	return ReviewOut{
		Name:        r.Name,
		Description: r.Description,
		ReviewURL:   r.ReviewURL,
		ReviewScore: r.ReviewScore,
		MediaType:   r.MediaType,
		GenreList:   r.GenreList,
		CreatorList: r.CreatorList,
	}

}
