package kiryuu

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"komikApi/controllers/komikCast"
	"net/http"
	"strings"
)

func fetchKomikInfoKiryuu(url string) (string, []map[string]string, string, []komikCast.Genre, string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", nil, "", nil, "", err
	}

	title := strings.TrimSpace(doc.Find(".entry-title").Text())

	sinopsis := strings.TrimSpace(doc.Find(".entry-content.entry-content-single p").Text())

	var chapters []map[string]string
	doc.Find(".clstyle li").Each(func(i int, s *goquery.Selection) {
		chapter := strings.TrimSpace(strings.Replace(s.Find(".chapternum").Text(), "Chapter ", "", -1))
		link, _ := s.Find(".eph-num a").Attr("href")
		time := strings.TrimSpace(s.Find(".chapterdate").Text())

		chapterInfo := map[string]string{
			"chapter": chapter,
			"link":    link,
			"time":    time,
		}
		chapters = append(chapters, chapterInfo)
	})

	status := strings.TrimSpace(doc.Find(".infotable tr:contains('Status') td:nth-child(2)").Text())
	typeInfo := strings.TrimSpace(doc.Find(".infotable tr:contains('Type') td:nth-child(2)").Text())
	var genres []komikCast.Genre
	genres = append(genres, komikCast.Genre{Name: status})
	genres = append(genres, komikCast.Genre{Name: "Kiryuu"})

	if typeInfo != "" {
		genres = append(genres, komikCast.Genre{Name: typeInfo})
	}

	doc.Find("div.seriestugenre a").Each(func(i int, s *goquery.Selection) {
		genreName := s.Text()
		genreLink, _ := s.Attr("href")

		genres = append(genres, komikCast.Genre{Name: genreName, Link: genreLink})
	})

	return title, chapters, sinopsis, genres, url, nil
}

func GetKomikInfoKiryuu(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	title, chapters, sinopsis, genres, link, err := fetchKomikInfoKiryuu(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":    title,
		"link":     link,
		"source":   "kiryuu",
		"sinopsis": sinopsis,
		"genre":    genres,
		"chapters": chapters,
	})
}
