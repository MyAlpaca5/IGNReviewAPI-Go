package schemas

import "time"

type Review struct {
	Name             string    `json:"name"`
	ShortName        string    `json:"short_name"`
	LongDescription  string    `json:"long_description"`
	ShortDescription string    `json:"short_description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ReviewURL        string    `json:"review_url"`
	ReviewScore      float32   `json:"review_score"`
	Slug             string    `json:"slug"`
	MediaType        string    `json:"media_type"`
	GenreList        []string  `json:"genre_list,omitempty"`
	CreatorList      []string  `json:"creator_list,omitempty"`
	PublisherList    []string  `json:"publisher_list,omitempty"`
	FranchiseList    []string  `json:"franchise_list,omitempty"`
	RegionList       []string  `json:"region_list,omitempty"`
}
