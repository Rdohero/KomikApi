package doujindesu

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"komikApi/allUrl"
	"komikApi/controllers/komikCast"
	"log"
	"net/http"
	"strings"
)

func fetchKomikInfoDoujindesu(url string) (string, []map[string]string, string, []komikCast.Genre, string, error) {

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
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return "", nil, "", nil, "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", nil, "", nil, "", err
	}

	title := strings.TrimSpace(strings.ReplaceAll(doc.Find("title").Text(), " - Doujindesu.XXX", ""))

	sinopsis := strings.TrimSpace(strings.Replace(doc.Find("div.pb-2 p").First().Text(), "Sinopsis:", "", -1))
	sinopsis = strings.ReplaceAll(sinopsis, "\\", "")
	sinopsis = strings.Replace(sinopsis, "Sinopsis :", "", -1)

	var chapters []map[string]string
	doc.Find("div.bxcl.scrolling ul li").Each(func(i int, s *goquery.Selection) {
		chapter := strings.TrimSpace(strings.Replace(s.Find("chapter").Text(), "Chapter ", "", -1))
		link, _ := s.Find(".epsright a").Attr("href")
		time := strings.TrimSpace(s.Find(".epsleft span.date").Text())

		chapterInfo := map[string]string{
			"chapter": chapter,
			"link":    allUrl.DoujindesuUrl + link,
			"time":    time,
		}
		chapters = append(chapters, chapterInfo)
	})

	status := strings.TrimSpace(doc.Find("tr:contains('Status') td a").Text())
	typeInfo := strings.TrimSpace(doc.Find("tr.magazines td a").Text())
	var genres []komikCast.Genre
	genres = append(genres, komikCast.Genre{Name: "doujindesu"})
	genres = append(genres, komikCast.Genre{Name: status})

	if typeInfo != "" {
		genres = append(genres, komikCast.Genre{Name: typeInfo})
	}

	doc.Find("div.tags a").Each(func(i int, s *goquery.Selection) {
		genreName := strings.TrimSpace(s.Text())
		genreLink, _ := s.Attr("href")

		genres = append(genres, komikCast.Genre{Name: genreName, Link: allUrl.DoujindesuUrl + genreLink})
	})

	return title, chapters, sinopsis, genres, url, nil
}

func GetKomikInfoDoujindesu(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	title, chapters, sinopsis, genres, link, err := fetchKomikInfoDoujindesu(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":    title,
		"link":     link,
		"source":   "doujindesu",
		"sinopsis": sinopsis,
		"genre":    genres,
		"chapters": chapters,
	})
}
