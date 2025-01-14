package pdf

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"os"
	"regexp"
)

func ExtractISSNNumbers(inputPath string) ([]string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	if numPages < 1 {
		return nil, fmt.Errorf("pdf tidak ada halamannya")
	}

	var issnNumbers []string

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return nil, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return nil, err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return nil, err
		}

		re := regexp.MustCompile(`(?i)ISSN\s*[:\-\s]*\D*(\d{4})[\s\-]*\D*(\d{4})`)
		matches := re.FindAllStringSubmatch(text, -1)

		if len(matches) > 0 {
			for _, match := range matches {
				issnNumber := match[1] + match[2]
				issnNumbers = append(issnNumbers, issnNumber)
			}
			break
		}
	}

	if len(issnNumbers) == 0 {
		return nil, nil
	}

	return issnNumbers, nil
}
