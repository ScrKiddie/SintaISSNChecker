package sinta

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func CheckSintaStatus(issn string) string {
	url := "https://sinta.kemdikbud.go.id/journals/?q=" + issn

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("error: status code %d\n", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	selection := doc.Find("span.num-stat.accredited a")
	if selection.Length() > 0 {
		accreditation := selection.Text()
		fmt.Println("status:", accreditation)
		if strings.Contains(accreditation, "Accredited") {
			return accreditation
		}
	}

	fmt.Println("bukan sinta")
	return ""
}
