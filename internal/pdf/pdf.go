package pdf

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

var pool pdfium.Pool
var instance pdfium.Pdfium
var regexPattern = `(?i)ISSN\s*[:\-\s]*\D*(\d{4})[\s\-]*\D*(\d{4})`
var initialized = false

func initPDFium() error {
	if initialized {
		return nil
	}

	originalStdout := os.Stdout
	originalStderr := os.Stderr

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("gagal membuka /dev/null: %w", err)
	}
	defer func() {
		if closeErr := devNull.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close /dev/null: %v\n", closeErr)
		}
	}()

	os.Stdout = devNull
	os.Stderr = devNull

	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	var initErr error
	pool, initErr = webassembly.Init(webassembly.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	})
	if initErr != nil {
		return fmt.Errorf("gagal menginisialisasi PDFium: %w", initErr)
	}

	instance, initErr = pool.GetInstance(time.Second * 30)
	if initErr != nil {
		return fmt.Errorf("gagal mendapatkan instance PDFium: %w", initErr)
	}

	initialized = true
	return nil
}

func ExtractISSNNumbers(ctx context.Context, inputPath string) ([]string, error) {
	if err := initPDFium(); err != nil {
		return nil, err
	}

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
			_ = err
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
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pageText, err := instance.GetPageText(&requests.GetPageText{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    pageNum,
				},
			},
		})
		if err != nil {
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
