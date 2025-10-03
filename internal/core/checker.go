package core

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

func CheckSintaStatus(issn string) (string, error) {
	url := "https://sinta.kemdikbud.go.id/journals/?q=" + issn
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := http.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("gagal melakukan HTTP GET request (percobaan %d/%d): %w", attempt, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		if res.StatusCode != http.StatusOK {
			_ = res.Body.Close()
			lastErr = fmt.Errorf("status code tidak 200, melainkan %d (percobaan %d/%d)", res.StatusCode, attempt, maxRetries)
			time.Sleep(retryDelay)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		_ = res.Body.Close()
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

	return "", fmt.Errorf("gagal memeriksa status SINTA setelah %d percobaan: %w", maxRetries, lastErr)
}
