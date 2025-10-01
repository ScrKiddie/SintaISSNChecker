package pdf

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

var pool pdfium.Pool
var instance pdfium.Pdfium
var logger *slog.Logger
var regexPattern = `(?i)ISSN\s*[:\-\s]*\D*(\d{4})[\s\-]*\D*(\d{4})`

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	var err error
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	})
	if err != nil {
		logger.Error("gagal menginisialisasi PDFium, program dihentikan", "error", err)
		os.Exit(1)
	}

	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		logger.Error("gagal mendapatkan instance PDFium, program dihentikan", "error", err)
		os.Exit(1)
	}
}

func ExtractISSNNumbers(inputPath string) ([]string, error) {
	pdfBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file PDF: %w", err)
	}

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka dokumen PDF: %w", err)
	}
	defer func(instance pdfium.Pdfium, request *requests.FPDF_CloseDocument) {
		_, err := instance.FPDF_CloseDocument(request)
		if err != nil {
			logger.Warn("gagal menutup dokumen PDF", "error", err)
		}
	}(instance, &requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan jumlah halaman: %w", err)
	}

	if pageCount.PageCount < 1 {
		return nil, fmt.Errorf("dokumen PDF tidak memiliki halaman")
	}

	var allText string

	for pageNum := 0; pageNum < pageCount.PageCount; pageNum++ {
		pageText, err := instance.GetPageText(&requests.GetPageText{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    pageNum,
				},
			},
		})
		if err != nil {
			logger.Warn("gagal mengekstrak teks dari halaman", "halaman", pageNum+1, "error", err)
			continue
		}

		allText += pageText.Text + "\n"

		re := regexp.MustCompile(regexPattern)
		matches := re.FindAllStringSubmatch(pageText.Text, -1)

		if len(matches) > 0 {
			var issnNumbers []string
			for _, match := range matches {
				issnNumber := match[1] + match[2]
				issnNumbers = append(issnNumbers, issnNumber)
			}
			return issnNumbers, nil
		}
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindAllStringSubmatch(allText, -1)

	if len(matches) == 0 {
		return nil, nil
	}

	var issnNumbers []string
	for _, match := range matches {
		issnNumber := match[1] + match[2]
		issnNumbers = append(issnNumbers, issnNumber)
	}

	return issnNumbers, nil
}
