package models

import (
	"errors"
	"regexp"
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

// Review struct includes all data related to one review record
type Review struct {
	dbBase
	ReviewOut
}

// ReviewOut struct includes data that is sent back via a response
type ReviewOut struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ReviewScore float32  `json:"review_score"`
	ReviewURL   string   `json:"review_url"`
	MediaType   string   `json:"media_type"`
	GenreList   []string `json:"genre_list"`
	CreatorList []string `json:"creator_list"`
}

func (r ReviewIn) ToReview() Review {
	var review Review

	review.Name = *r.Name
	review.ReviewScore = *r.ReviewScore
	review.ReviewURL = *r.ReviewURL

	if r.Description != nil {
		review.Description = *r.Description
	}

	if r.MediaType != nil {
		review.MediaType = *r.MediaType
	}

	if r.GenreList != nil {
		review.GenreList = r.GenreList
	}

	if r.CreatorList != nil {
		review.CreatorList = r.CreatorList
	}

	return review
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

func (r Review) ToReviewOut() ReviewOut {
	return r.ReviewOut
}
