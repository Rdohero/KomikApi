package kiryuu

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func fetchDataFromURLKiryuu(pageURL string) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitVisible(`div#readerarea`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, fmt.Errorf("error navigating to the page: %v", err)
	}

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
