package main

import (
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sinta/internal/logging"
	"sinta/internal/ui"
	"syscall"

	"github.com/gonutz/w32/v2"
	"github.com/ncruces/zenity"
)

func main() {
	hideConsoleWindow()

	if !zenity.IsAvailable() {
		slog.Error("Dialog GUI tidak tersedia. Silakan install zenity terlebih dahulu.")
		showErrorMessage("Dialog GUI tidak tersedia. Silakan install zenity terlebih dahulu.")
		os.Exit(1)
	}

	logHandler := logging.NewCaptureHandler(slog.NewTextHandler(os.Stdout, nil))
	logger := slog.New(logHandler)

	defer func() {
		cleanupTempFiles()
		logger.Info("Aplikasi selesai, file temporary dibersihkan.")
	}()

	setupSignalHandler(logger)

	appUI := ui.NewAppUI(logger, logHandler)
	appUI.Run()
}

func hideConsoleWindow() {
	if runtime.GOOS != "windows" {
		return
	}

	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, w32.SW_HIDE)
		}
	}
}

func setupSignalHandler(logger *slog.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-sigChan
		cleanupTempFiles()
		logger.Info("Aplikasi dihentikan, file temporary dibersihkan.")
		os.Exit(0)
	}()
}

func cleanupTempFiles() {
	tempDir := os.TempDir()

	patterns := []string{
		"sinta_log.txt",
		"sinta_temp_*.txt",
		"zenity_*",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(tempDir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			_ = os.Remove(match)
		}
	}
}

func showErrorMessage(message string) {
	errorFile := filepath.Join(os.TempDir(), "sinta-error.txt")
	_ = os.WriteFile(errorFile, []byte(message), 0644)

	if runtime.GOOS == "windows" {
		cmd := exec.Command("msg", "*", message)
		_ = cmd.Run()
	}
}
