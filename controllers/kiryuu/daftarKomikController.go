package kiryuu

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"komikApi/allUrl"
	"net/http"
	"strconv"
	"strings"
)

type DaftarKomikModelKiryuu struct {
	Title    string `json:"title"`
	Thumb    string `json:"thumb"`
	Chapter  string `json:"chapter"`
	Rating   string `json:"rating"`
	Source   string `json:"source"`
	KomikURL string `json:"komikUrl"`
}

type KomikResponseKiryuu struct {
	DaftarKomik    []DaftarKomikModelKiryuu `json:"daftar_komik"`
	PaginationPage int                      `json:"page"`
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
			Source:   "kiryuu",
			Chapter:  chapter,
			Rating:   strings.TrimSpace(rating),
			KomikURL: link,
		})
	})

	return KomikResponseKiryuu{
		DaftarKomik:    daftarKomik,
		PaginationPage: paginationPage,
	}, nil
}
