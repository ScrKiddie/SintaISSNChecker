package sinta

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CheckSintaStatus(issn string) (string, error) {
	url := "https://sinta.kemdikbud.go.id/journals/?q=" + issn

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("gagal melakukan HTTP GET request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code tidak 200, melainkan: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("gagal mem-parsing response body: %w", err)
	}

	selection := doc.Find("span.num-stat.accredited a")
	if selection.Length() > 0 {
		accreditation := selection.Text()
		if strings.Contains(accreditation, "Accredited") {
			return accreditation, nil
		}
	}
	return "", nil
}
