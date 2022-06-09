package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MyAlpaca5/IGNReviewAPI-Go/internal/api/schemas"
	"github.com/gin-gonic/gin"
)

// ReviewsGETHandler handles "POST /api/reviews" endpoint. TODO: for now, just return plain text.
func CreateReviewHandler(c *gin.Context) {
	// TEST
	c.String(http.StatusOK, "POST %s", c.FullPath())
	// TEST END
}

// ReviewsGETHandler handles "GET /api/reviews/:id" endpoint. TODO: for now, just return plain text.
func ShowReviewHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(id)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid id parameter: %s", c.Param("id"))})
		return
	}

	// TEST
	r := schemas.Review{
		Name:             strconv.Itoa(int(id)),
		ShortName:        "short test",
		LongDescription:  "",
		ShortDescription: "short test",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now().Add(time.Hour),
		ReviewURL:        "url to here",
		ReviewScore:      4.5,
		Slug:             "test me",
		MediaType:        "abc",
		GenreList:        nil,
		CreatorList:      []string{},
		PublisherList:    []string{"a"},
		FranchiseList:    []string{"a", "b"},
		RegionList:       []string{"a", "b", "c"},
	}

	c.JSON(http.StatusOK, r)
	// TEST END
}
