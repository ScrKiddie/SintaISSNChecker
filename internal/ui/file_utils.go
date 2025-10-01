package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validatePDFFolder(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("gagal akses folder: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path yang dipilih bukan folder")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca isi folder: %w", err)
	}

	var pdfFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.ToLower(filepath.Ext(entry.Name())) == ".pdf" {
			pdfFiles = append(pdfFiles, filepath.Join(path, entry.Name()))
		}
	}

	if len(pdfFiles) == 0 {
		return nil, fmt.Errorf("tidak ada file PDF yang ditemukan di dalam folder yang dipilih")
	}

	return pdfFiles, nil
}
