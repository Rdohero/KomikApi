package komikCast

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"net/http"
	"strings"
)

type Genre struct {
	Name string
	Link string
}

func fetchGenres(url string) ([]Genre, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var genres []Genre
	doc.Find(".genresx li a").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		link, exists := s.Attr("href")
		if exists {
			genres = append(genres, Genre{Name: name, Link: link})
		}
	})

	return genres, nil
}

func GetGenreInfo(c *gin.Context) {
	url := allUrl.KomikCastUrl + "/genres/fantasy/"
	genres, err := fetchGenres(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}
