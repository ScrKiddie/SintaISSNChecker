package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sinta/internal/ui"
	"sinta/internal/utility"
	"syscall"

	"github.com/ncruces/zenity"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("Application crashed: %v", r)
			showErrorMessage(errorMsg)
			os.Exit(1)
		}
	}()

	if !zenity.IsAvailable() {
		showErrorMessage("Dialog GUI tidak tersedia. Silakan install zenity terlebih dahulu.")
		os.Exit(1)
	}

	discardHandler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logHandler := utility.NewCaptureHandler(discardHandler)
	logger := slog.New(logHandler)

	defer func() {
		cleanupTempFiles()
	}()

	setupSignalHandler()

	appUI := ui.NewAppUI(logger, logHandler)
	appUI.Run()
}

func setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-sigChan
		cleanupTempFiles()
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
	if runtime.GOOS == "windows" {
		methods := []func(string) error{
			func(msg string) error {
				return exec.Command("msg", "*", msg).Run()
			},
			func(msg string) error {
				return exec.Command("powershell", "-Command", fmt.Sprintf("[System.Windows.Forms.MessageBox]::Show('%s', 'Error', 'OK', 'Error')", msg)).Run()
			},
			func(msg string) error {
				errorFile := filepath.Join(os.TempDir(), "sinta-error.txt")
				return os.WriteFile(errorFile, []byte(msg), 0644)
			},
		}

		for _, method := range methods {
			if err := method(message); err == nil {
				break
			}
		}
	}
}
