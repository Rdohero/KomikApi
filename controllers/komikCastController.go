package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Komik struct {
	Title         string `json:"title"`
	Chapter       string `json:"chapter"`
	Rating        string `json:"rating"`
	Image         string `json:"image"`
	Type          string `json:"type"`
	IsCompleted   bool   `json:"isCompleted"`
	Link          string `json:"link"`
	LinkId        string `json:"linkId"`
	LinkChapter   string `json:"linkChapter"`
	LinkIdChapter string `json:"linkIdChapter"`
}

type KomikResponse struct {
	DaftarKomik    []Komik `json:"daftarKomik"`
	PaginationPage int     `json:"page"`
}

func GetDaftarKomik(order string, page string) (KomikResponse, error) {
	urlPath := allUrl.KomikCastUrl + "/daftar-komik/"

	if order != "" {
		order = "?order=" + order
	}

	if page != "" {
		urlPath += "page/" + page + "/" + order
	} else {
		urlPath += order
	}

	paginationPage, err := strconv.Atoi(page)
	if err != nil {
		paginationPage = 1
	}

	resp, err := http.Get(urlPath)
	if err != nil {
		return KomikResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return KomikResponse{}, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return KomikResponse{}, err
	}

	var daftarKomik []Komik

	doc.Find("div.list-update_item").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")
		linkId := strings.TrimPrefix(link, allUrl.KomikCastUrl+"/komik/")
		linkId = strings.TrimSuffix(linkId, "/")

		title := s.Find("h3.title").Text()
		chapter := strings.TrimSpace(strings.Replace(s.Find("div.chapter").Text(), "Ch.", "", 1))
		rating := s.Find("div.numscore").Text()
		image, _ := s.Find("img").Attr("src")
		komikType := s.Find("span.type").Text()
		isCompleted := s.Find("span.Completed").Length() > 0

		linkChapter, _ := s.Find("div.chapter").Attr("href")
		linkIdChapter := strings.TrimPrefix(linkChapter, allUrl.KomikCastUrl+"/chapter/")
		linkIdChapter = strings.TrimSuffix(linkIdChapter, "/")

		daftarKomik = append(daftarKomik, Komik{
			Title:         strings.TrimSpace(title),
			Chapter:       chapter,
			Rating:        strings.TrimSpace(rating),
			Image:         strings.TrimSpace(image),
			Type:          strings.TrimSpace(komikType),
			IsCompleted:   isCompleted,
			Link:          link,
			LinkId:        linkId,
			LinkChapter:   linkChapter,
			LinkIdChapter: linkIdChapter,
		})
	})

	return KomikResponse{
		DaftarKomik:    daftarKomik,
		PaginationPage: paginationPage,
	}, nil
}

func fetchDataFromURL(url string) (string, []string, error) {
	// Send HTTP request to the URL
	res, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", nil, err
	}

	// Extract text from &lt;h1&gt; element with itemprop="name"
	title := doc.Find("h1[itemprop='name']").Text()

	// Extract image URLs from &lt;div class="main-reading-area"&gt;
	var imgUrls []string
	doc.Find("div.main-reading-area img").Each(func(i int, s *goquery.Selection) {
		imgUrl, exists := s.Attr("src")
		if exists {
			imgUrls = append(imgUrls, imgUrl)
		}
	})

	return title, imgUrls, nil
}

// fetchChapterURLs fetches the next and previous chapter URLs from the specified URL
func fetchChapterURLs(url string) (string, string, error) {
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
	doc.Find("div.nextprev a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "Next Chapter") {
			nextChapterURL, _ = s.Attr("href")
		}
		if strings.Contains(text, "Previous Chapter") {
			prevChapterURL, _ = s.Attr("href")
		}
	})

	if nextChapterURL == "" && prevChapterURL == "" {
		return "", "", errors.New("chapter URLs not found")
	}

	return nextChapterURL, prevChapterURL, nil
}

// GetDataHandler handles the request and returns the title, image URLs, and chapter URLs
func GetDataHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	title, imgUrls, err := fetchDataFromURL(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	nextChapterURL, prevChapterURL, err := fetchChapterURLs(url)
	if err != nil {
		fmt.Println("Warning: Chapter URLs not found:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"title":          title,
		"imgUrls":        imgUrls,
		"nextChapterURL": nextChapterURL,
		"prevChapterURL": prevChapterURL,
	})
}

type Genres struct {
	Name string
	Link string
}

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
		"sinopsis": sinopsis,
		"genre":    genres,
		"chapters": chapters,
	})
}

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
