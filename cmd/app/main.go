package main

import (
	"fmt"
	"log/slog"
	"os"
	"sinta/internal/pdf"
	"sinta/internal/sinta"
	"strings"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	files, err := os.ReadDir("./pdf")
	if err != nil {
		logger.Error("gagal membaca direktori", "error", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		logger.Warn("tidak ada file di dalam direktori 'pdf', program dihentikan")
		return
	}

	for i, file := range files {
		if strings.HasSuffix(file.Name(), ".pdf") {
			filePath := "./pdf/" + file.Name()
			logger.Info("memulai pemrosesan file", "file_number", i+1, "file_path", filePath)

			issnNumbers, err := pdf.ExtractISSNNumbers(filePath)
			if err != nil {
				logger.Error("gagal mengekstrak ISSN dari file", "file_path", filePath, "error", err)
				continue
			}

			if len(issnNumbers) == 0 {
				logger.Warn("tidak ditemukan ISSN di dalam file", "file_path", filePath)
				continue
			}

			hasSinta := false
			var accreditation string
			for _, issn := range issnNumbers {
				logger.Info("memeriksa status SINTA untuk ISSN", "issn", issn)
				accreditation, err = sinta.CheckSintaStatus(issn)
				if err != nil {
					logger.Error("gagal memeriksa status SINTA", "issn", issn, "error", err)
					continue
				}
				if accreditation != "" {
					hasSinta = true
					logger.Info("ditemukan akreditasi SINTA", "issn", issn, "akreditasi", accreditation)
					break
				}
			}

			if hasSinta {
				expectedPrefix := fmt.Sprintf("%s - ", strings.TrimSpace(accreditation))
				if strings.HasPrefix(file.Name(), expectedPrefix) {
					logger.Info("file sudah memiliki akreditasi SINTA di namanya, tidak diubah", "file_path", filePath, "akreditasi", accreditation)
					continue
				}

				newFileName := fmt.Sprintf("./pdf/%s - %s.pdf", strings.TrimSpace(accreditation), strings.TrimSuffix(file.Name(), ".pdf"))
				err := os.Rename(filePath, newFileName)
				if err != nil {
					logger.Error("gagal mengganti nama file", "old_path", filePath, "new_path", newFileName, "error", err)
				} else {
					logger.Info("berhasil mengganti nama file", "new_path", newFileName)
				}
			} else {
				logger.Info("jurnal tidak terakreditasi SINTA, nama file tidak diubah", "file_path", filePath)
			}
		}
	}
}
