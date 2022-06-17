package models

// Review struct includes all data related to one review record
type Review struct {
	dbBase
	Name        string
	Description string
	ReviewURL   string
	ReviewScore float32
	MediaType   string
	GenreList   []string
	CreatorList []string
}
