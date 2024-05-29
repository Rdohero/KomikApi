package komikCast

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"komikApi/allUrl"
	"net/http"
	"strconv"
	"strings"
)

type Komik struct {
	Title   string `json:"title"`
	Chapter string `json:"chapter"`
	Rating  string `json:"rating"`
	Source  string `json:"source"`
	Image   string `json:"image"`
	Link    string `json:"link"`
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

		linkChapter, _ := s.Find("div.chapter").Attr("href")
		linkIdChapter := strings.TrimPrefix(linkChapter, allUrl.KomikCastUrl+"/chapter/")
		linkIdChapter = strings.TrimSuffix(linkIdChapter, "/")

		daftarKomik = append(daftarKomik, Komik{
			Title:   strings.TrimSpace(title),
			Chapter: chapter,
			Rating:  strings.TrimSpace(rating),
			Source:  "komik_cast",
			Image:   strings.TrimSpace(image),
			Link:    link,
		})
	})

	return KomikResponse{
		DaftarKomik:    daftarKomik,
		PaginationPage: paginationPage,
	}, nil
}
