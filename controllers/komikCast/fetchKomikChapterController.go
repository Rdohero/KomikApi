package komikCast

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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
