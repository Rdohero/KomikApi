package doujindesu

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"log"
	"net/http"
	"strings"
)

func fetchDataFromURLDoujindesu(pageURL string) ([]string, string, string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("incognito", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageURL),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, "", "", fmt.Errorf("error navigating to the page: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, "", "", fmt.Errorf("error navigating to the page: %v", err)
	}

	var imgUrls []string
	doc.Find("div.main div div#anu img").Each(func(i int, s *goquery.Selection) {
		imgUrl, exists := s.Attr("src")
		if exists {
			imgUrls = append(imgUrls, imgUrl)
		}
	})

	var nextChapterURL, prevChapterURL string
	prevChapterURL = doc.Find("div.nvs a").AttrOr("href", "")
	nextChapterURL = doc.Find("div.nvs.rght a").AttrOr("href", "")

	return imgUrls, nextChapterURL, prevChapterURL, nil
}

func GetDataHandlerDoujindesu(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	imgUrls, nextChapterURL, prevChapterURL, err := fetchDataFromURLDoujindesu(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imgUrls":        imgUrls,
		"nextChapterURL": allUrl.DoujindesuUrl + nextChapterURL,
		"prevChapterURL": allUrl.DoujindesuUrl + prevChapterURL,
	})
}
