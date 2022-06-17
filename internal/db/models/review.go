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

// ReviewQueryParam struct includes accepted filter for read action
type ReviewQueryParam struct {
	Name     *string
	ScoreMin *string
	Order    *string
	Page     *string
	PageSize *string
	Genres   []string
}

func NewReviewQueryParam(queryParam map[string][]string) ReviewQueryParam {
	var reviewQueryParam ReviewQueryParam
	if name, found := queryParam["name"]; found {
		reviewQueryParam.Name = &(name[0])
	}
	if scoreMin, found := queryParam["score_min"]; found {
		reviewQueryParam.ScoreMin = &(scoreMin[0])
	}
	if order, found := queryParam["order"]; found {
		reviewQueryParam.Order = &(order[0])
	}
	if page, found := queryParam["page"]; found {
		reviewQueryParam.Page = &(page[0])
	}
	if pageSize, found := queryParam["page_size"]; found {
		reviewQueryParam.PageSize = &(pageSize[0])
	}
	if genres, found := queryParam["genres"]; found {
		reviewQueryParam.Genres = genres
	}
	return reviewQueryParam
}
