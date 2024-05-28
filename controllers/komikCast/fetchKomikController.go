package komikCast

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func fetchKomikInfo(url string) (string, []map[string]string, string, []Genre, string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", nil, "", nil, "", err
	}

	title := strings.TrimSpace(doc.Find(".komik_info-content-body-title").Text())

	sinopsis1 := strings.TrimSpace(doc.Find(".komik_info-description-sinopsis p").Text())

	sinopsis2 := strings.TrimSpace(doc.Find(".markup-2BOw-j.messageContent-2qWWxC").Text())

	var sinopsis string

	if sinopsis1 != "" {
		sinopsis = sinopsis1
	} else if sinopsis2 != "" {
		sinopsis = sinopsis2
	} else {
		sinopsis = "Sinopsis tidak ditemukan"
	}

	var chapters []map[string]string
	doc.Find(".komik_info-chapters-item").Each(func(i int, s *goquery.Selection) {
		chapter := strings.TrimSpace(strings.Replace(s.Find(".chapter-link-item").Text(), "Chapter\n", "", -1))
		link, _ := s.Find(".chapter-link-item").Attr("href")
		time := strings.TrimSpace(s.Find(".chapter-link-time").Text())

		chapterInfo := map[string]string{
			"chapter": chapter,
			"link":    link,
			"time":    time,
		}
		chapters = append(chapters, chapterInfo)
	})

	status := strings.TrimSpace(doc.Find(".komik_info-content-info b:contains('Status:')").Parent().Contents().Last().Text())
	typeInfo := strings.TrimSpace(doc.Find(".komik_info-content-info-type a").Text())
	typeInfoLink, _ := doc.Find(".komik_info-content-info-type a").Attr("href")

	var genres []Genre
	genres = append(genres, Genre{Name: status})
	genres = append(genres, Genre{Name: "Komik Cast"})

	if typeInfo != "" && typeInfoLink != "" {
		genres = append(genres, Genre{Name: typeInfo, Link: typeInfoLink})
	}

	doc.Find("span.komik_info-content-genre a.genre-item").Each(func(i int, s *goquery.Selection) {
		genreName := strings.TrimSpace(s.Text())
		genreLink, _ := s.Attr("href")
		genres = append(genres, Genre{Name: genreName, Link: genreLink})
	})

	return title, chapters, sinopsis, genres, url, nil
}

func GetKomikInfo(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	title, chapters, sinopsis, genres, link, err := fetchKomikInfo(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":    title,
		"link":     link,
		"source":   "komik_cast",
		"sinopsis": sinopsis,
		"genre":    genres,
		"chapters": chapters,
	})
}
