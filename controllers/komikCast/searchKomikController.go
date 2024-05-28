package komikCast

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"net/http"
	"strings"
)

type KomikSearch struct {
	Title      string `json:"title"`
	Thumb      string `json:"thumb"`
	Type       string `json:"type"`
	Chapter    string `json:"chapter"`
	Rating     string `json:"rating"`
	KomikURL   string `json:"komikUrl"`
	ChapterURL string `json:"chapterUrl"`
}

func SearchKomik(c *gin.Context) {
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")

	search = strings.ReplaceAll(search, " ", "+")
	baseURL := fmt.Sprintf(allUrl.KomikCastUrl+"/page/%s/?s=%s", page, search)

	resp, err := http.Get(baseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": resp.Status})
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var listOfKomik []KomikSearch
	doc.Find(".list-update_item").Each(func(_ int, s *goquery.Selection) {
		komikURL, _ := s.Find("a.data-tooltip").Attr("href")
		thumb, _ := s.Find(".list-update_item-image img").Attr("src")
		title := s.Find(".list-update_item-info .title").Text()
		chapter := s.Find(".list-update_item-info .chapter").Text()
		rating := s.Find(".numscore").Text()

		k := KomikSearch{
			Title:      strings.TrimSpace(title),
			Thumb:      thumb,
			Type:       s.Find(".list-update_item-image .type").Text(),
			Chapter:    strings.TrimSpace(chapter),
			Rating:     rating,
			KomikURL:   komikURL,
			ChapterURL: s.Find(".list-update_item-info .chapter").AttrOr("href", ""),
		}
		listOfKomik = append(listOfKomik, k)
	})

	if len(listOfKomik) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Judul Yang dicari tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, listOfKomik)
}
