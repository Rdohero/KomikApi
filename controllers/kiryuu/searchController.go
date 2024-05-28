package kiryuu

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"net/http"
	"strings"
)

type KomikSearchKiryuu struct {
	Title    string `json:"title"`
	Thumb    string `json:"thumb"`
	Source   string `json:"source"`
	Chapter  string `json:"chapter"`
	Rating   string `json:"rating"`
	KomikURL string `json:"komikUrl"`
}

func SearchKomikKiryuu(c *gin.Context) {
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")

	search = strings.ReplaceAll(search, " ", "+")
	baseURL := fmt.Sprintf(allUrl.KiryuuUrl+"/page/%s/?s=%s", page, search)

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

	var listOfKomik []KomikSearchKiryuu
	doc.Find(".bs").Each(func(_ int, s *goquery.Selection) {
		komikURL, _ := s.Find("a").Attr("href")
		thumb, _ := s.Find(".limit img").Attr("src")
		title := s.Find(".tt").Text()
		chapter := s.Find(".epxs").Text()
		rating := s.Find(".numscore").Text()

		k := KomikSearchKiryuu{
			Title:    strings.TrimSpace(title),
			Source:   "kiryuu",
			Thumb:    thumb,
			Chapter:  strings.TrimSpace(chapter),
			Rating:   rating,
			KomikURL: komikURL,
		}
		listOfKomik = append(listOfKomik, k)
	})

	if len(listOfKomik) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Judul Yang dicari tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, listOfKomik)
}
