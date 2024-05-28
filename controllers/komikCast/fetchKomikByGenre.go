package komikCast

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Comic struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	ImageLink string `json:"image_link"`
}

func FetchComicsByGenre(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comics"})
		return
	}

	var comics []Comic

	doc.Find(".list-update_item").Each(func(i int, s *goquery.Selection) {
		comic := Comic{}

		// Get title
		comic.Title = s.Find(".title").Text()

		// Get link
		link, exists := s.Find("a").Attr("href")
		if exists {
			comic.Link = link
		}

		// Get image link
		imageLink, exists := s.Find("img").Attr("src")
		if exists {
			comic.ImageLink = imageLink
		}

		comics = append(comics, comic)
	})

	c.JSON(http.StatusOK, comics)
}
