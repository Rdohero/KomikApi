package controllers

import (
	"errors"
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

func fetchKomikInfoKiryuu(url string) (string, []map[string]string, string, []Genre, string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", nil, "", nil, "", err
	}

	title := strings.TrimSpace(doc.Find(".entry-title").Text())

	sinopsis := strings.TrimSpace(doc.Find(".entry-content.entry-content-single p").Text())

	var chapters []map[string]string
	doc.Find(".clstyle li").Each(func(i int, s *goquery.Selection) {
		chapter := strings.TrimSpace(strings.Replace(s.Find(".chapternum").Text(), "Chapter\n", "", -1))
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
	var genres []Genre
	genres = append(genres, Genre{Name: status})

	if typeInfo != "" {
		genres = append(genres, Genre{Name: typeInfo})
	}

	doc.Find("div.seriestugenre a").Each(func(i int, s *goquery.Selection) {
		genreName := s.Text()
		genreLink, _ := s.Attr("href")

		genres = append(genres, Genre{Name: genreName, Link: genreLink})
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
		"sinopsis": sinopsis,
		"genre":    genres,
		"chapters": chapters,
	})
}

func fetchDataFromURLKiryuu(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var imgUrls []string
	doc.Find("div#readerarea img.ts-main-image").Each(func(i int, s *goquery.Selection) {
		imgUrl, exists := s.Attr("src")
		if exists {
			imgUrls = append(imgUrls, imgUrl)
		}
	})

	return imgUrls, nil
}

func fetchChapterURLsKiryuu(url string) (string, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	var nextChapterURL, prevChapterURL string
	prevChapterRel := doc.Find("div.nextprev a.ch-prev-btn").AttrOr("href", "")
	nextChapterRel := doc.Find("div.nextprev a.ch-next-btn").AttrOr("href", "")
	prevChapterURL = url + prevChapterRel
	nextChapterURL = url + nextChapterRel
	if nextChapterURL == "" && prevChapterURL == "" {
		return "", "", errors.New("chapter URLs not found")
	}

	return nextChapterURL, prevChapterURL, nil
}

func GetDataHandlerKiryuu(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	imgUrls, err := fetchDataFromURLKiryuu(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	nextChapterURL, prevChapterURL, err := fetchChapterURLsKiryuu(url)
	if err != nil {
		fmt.Println("Warning: Chapter URLs not found:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"imgUrls":        imgUrls,
		"nextChapterURL": nextChapterURL,
		"prevChapterURL": prevChapterURL,
	})
}
