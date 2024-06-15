package doujindesu

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"komikApi/allUrl"
	"log"
	"strings"
)

type DaftarKomikModelDoujindesu struct {
	Title   string `json:"title"`
	Chapter string `json:"chapter"`
	Rating  string `json:"rating"`
	Source  string `json:"source"`
	Image   string `json:"image"`
	Link    string `json:"link"`
}

type KomikResponseDoujindesu struct {
	DaftarKomik    []DaftarKomikModelDoujindesu `json:"daftarKomik"`
	PaginationPage int                          `json:"page"`
}

func GetDaftarKomikDoujindesu() (KomikResponseDoujindesu, error) {
	urlPath := allUrl.DoujindesuUrl

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
		chromedp.Navigate(urlPath),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return KomikResponseDoujindesu{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return KomikResponseDoujindesu{}, err
	}

	var daftarKomik []DaftarKomikModelDoujindesu

	doc.Find("div.entries article.entry").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")

		title, _ := s.Find("a").Attr("title")
		cText := strings.TrimSpace(s.Find("div.metadata div.artists a span").Text())
		chapter := strings.TrimSpace(strings.Replace(cText, "Chapter", "", 1))
		chapter = strings.TrimSpace(strings.Replace(chapter, "END", "", 1))
		image, _ := s.Find("img").Attr("src")

		daftarKomik = append(daftarKomik, DaftarKomikModelDoujindesu{
			Title:   strings.TrimSpace(title),
			Image:   image,
			Source:  "doujindesu",
			Chapter: chapter,
			Link:    urlPath + link,
		})
	})

	return KomikResponseDoujindesu{
		DaftarKomik: daftarKomik,
	}, err
}
