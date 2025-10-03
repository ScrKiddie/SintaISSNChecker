package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sinta/internal/pdf"
	"strings"
)

type ProcessResult struct {
	SuccessFiles         []string
	FailedFiles          []string
	NotAccreditedFiles   []string
	AlreadyRenamedFiles  []string
	ProcessingErrorFiles []string
	ISSNNotFoundFiles    []string
}

type ProgressUpdater func(percent int, text string) error

func ProcessFiles(ctx context.Context, logger *slog.Logger, files []string, updateProgress ProgressUpdater) ProcessResult {
	logger.Info("Memulai pemrosesan file PDF", "count", len(files))

	var result ProcessResult

	for i, filePath := range files {
		percent := (i * 100) / len(files)
		if err := updateProgress(percent, fmt.Sprintf("Memproses %d/%d: %s", i+1, len(files), filepath.Base(filePath))); err != nil {
			logger.Warn("Progress dialog ditutup, pemrosesan dihentikan.", "error", err)
			break
		}

		logger.Info("Memproses file", "file", filepath.Base(filePath), "current", i+1, "total", len(files))

		select {
		case <-ctx.Done():
			logger.Info("Pemrosesan dibatalkan oleh pengguna.")
			return result
		default:
		}

		issnNumbers, err := pdf.ExtractISSNNumbers(ctx, filePath)
		if err != nil {
			logger.Error("Gagal mengekstrak ISSN", "file", filepath.Base(filePath), "error", err.Error())
			result.ProcessingErrorFiles = append(result.ProcessingErrorFiles, filePath)
			continue
		}

		if len(issnNumbers) == 0 {
			logger.Warn("Tidak ditemukan ISSN di dalam file", "file", filepath.Base(filePath))
			result.ISSNNotFoundFiles = append(result.ISSNNotFoundFiles, filePath)
			continue
		}

		logger.Info("ISSN ditemukan", "file", filepath.Base(filePath), "issns", strings.Join(issnNumbers, ", "))

		hasSinta := false
		var accreditation string
		var sintaCheckFailed bool

		for _, issn := range issnNumbers {
			logger.Info("Memeriksa status SINTA", "issn", issn)
			accreditation, err = CheckSintaStatus(issn)
			if err != nil {
				logger.Error("Gagal memeriksa status SINTA", "issn", issn, "error", err.Error())
				sintaCheckFailed = true
				continue
			}
			if accreditation != "" {
				hasSinta = true
				logger.Info("Akreditasi SINTA ditemukan", "issn", issn, "accreditation", accreditation)
				break
			}
		}

		if sintaCheckFailed && !hasSinta {
			logger.Error("Gagal memeriksa status SINTA untuk semua ISSN", "file", filepath.Base(filePath))
			result.ProcessingErrorFiles = append(result.ProcessingErrorFiles, filePath)
			continue
		}

		if hasSinta {
			fileName := filepath.Base(filePath)
			expectedPrefix := fmt.Sprintf("%s - ", strings.TrimSpace(accreditation))
			if strings.HasPrefix(fileName, expectedPrefix) {
				logger.Info("File sudah memiliki akreditasi SINTA di namanya", "file", filePath)
				result.AlreadyRenamedFiles = append(result.AlreadyRenamedFiles, filePath)
				continue
			}

			dir := filepath.Dir(filePath)
			newFileName := fmt.Sprintf("%s/%s - %s.pdf", dir, strings.TrimSpace(accreditation), strings.TrimSuffix(fileName, ".pdf"))
			err := os.Rename(filePath, newFileName)
			if err != nil {
				logger.Error("Gagal mengganti nama file", "file", filepath.Base(filePath), "error", err.Error())
				result.FailedFiles = append(result.FailedFiles, filePath)
			} else {
				logger.Info("Berhasil mengganti nama file", "old", filepath.Base(filePath), "new", filepath.Base(newFileName))
				result.SuccessFiles = append(result.SuccessFiles, newFileName)
			}
		} else {
			logger.Info("Jurnal tidak terakreditasi SINTA", "file", filepath.Base(filePath))
			result.NotAccreditedFiles = append(result.NotAccreditedFiles, filePath)
		}
	}

	return result
}
