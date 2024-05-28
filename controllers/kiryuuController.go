package controllers

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"net/http"
	"strconv"
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

func fetchDataFromURLKiryuu(pageURL string) ([]string, error) {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Navigate to the URL and retrieve the HTML content
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitVisible(`div#readerarea`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, fmt.Errorf("error navigating to the page: %v", err)
	}

	// Process the HTML content using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
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

func fetchChapterURLsKiryuu(pageURL string) (string, string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitVisible(`div.nextprev`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return "", "", fmt.Errorf("error navigating to the page: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", "", fmt.Errorf("error parsing HTML: %v", err)
	}

	var nextChapterURL, prevChapterURL string
	prevChapterURL = doc.Find("div.nextprev a.ch-prev-btn").AttrOr("href", "")
	nextChapterURL = doc.Find("div.nextprev a.ch-next-btn").AttrOr("href", "")

	if nextChapterURL == "" && prevChapterURL == "" {
		return "", "", fmt.Errorf("chapter URLs not found")
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

type DaftarKomikModelKiryuu struct {
	Title    string `json:"title"`
	Thumb    string `json:"thumb"`
	Chapter  string `json:"chapter"`
	Rating   string `json:"rating"`
	KomikURL string `json:"komikUrl"`
}

type KomikResponseKiryuu struct {
	DaftarKomikModelKiryuu []DaftarKomikModelKiryuu `json:"daftar_komik_model_kiryuu"`
	PaginationPage         int                      `json:"page"`
}

func GetDaftarKomikKiryuu(page string) (KomikResponseKiryuu, error) {
	urlPath := allUrl.KiryuuUrl + "/manga/"

	if page != "" {
		urlPath += "?page=" + page
	}

	paginationPage, err := strconv.Atoi(page)
	if err != nil {
		paginationPage = 1
	}

	resp, err := http.Get(urlPath)
	if err != nil {
		return KomikResponseKiryuu{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return KomikResponseKiryuu{}, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return KomikResponseKiryuu{}, err
	}

	var daftarKomik []DaftarKomikModelKiryuu

	doc.Find("div.listupd div.bs").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")

		title := s.Find("div.tt").Text()
		chapter := strings.TrimSpace(strings.Replace(s.Find("div.epxs").Text(), "Chapter", "", 1))
		rating := s.Find("div.numscore").Text()
		image, _ := s.Find("img.ts-post-image").Attr("src")

		daftarKomik = append(daftarKomik, DaftarKomikModelKiryuu{
			Title:    strings.TrimSpace(title),
			Thumb:    image,
			Chapter:  chapter,
			Rating:   strings.TrimSpace(rating),
			KomikURL: link,
		})
	})

	return KomikResponseKiryuu{
		DaftarKomikModelKiryuu: daftarKomik,
		PaginationPage:         paginationPage,
	}, nil
}
