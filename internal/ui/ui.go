package ui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sinta/internal/core"
	"sinta/internal/utility"
	"strings"

	"github.com/ncruces/zenity"
)

type AppUI struct {
	logger  *slog.Logger
	logHndl *utility.CaptureHandler
}

func NewAppUI(logger *slog.Logger, logHndl *utility.CaptureHandler) *AppUI {
	return &AppUI{
		logger:  logger,
		logHndl: logHndl,
	}
}

func (app *AppUI) Run() {
	for {
		choice, err := zenity.List(
			"Pilih Aksi:",
			[]string{
				"Pilih Folder PDF",
				"Pilih File PDF (Multiple)",
				"Lihat Log",
			},
			zenity.Title("SINTA ISSN Checker by ScrKiddie"),
			zenity.Width(500),
			zenity.Height(350),
		)

		if err != nil {
			if errors.Is(err, zenity.ErrCanceled) {
				return
			}
			continue
		}

		switch choice {
		case "Pilih Folder PDF":
			app.selectFolder()
		case "Pilih File PDF (Multiple)":
			app.selectFiles()
		case "Lihat Log":
			app.showLogs()
		}
	}
}

func (app *AppUI) selectFolder() {
	folderPath, err := zenity.SelectFile(
		zenity.Title("Pilih Folder PDF"),
		zenity.Directory(),
	)
	if err != nil {
		if !errors.Is(err, zenity.ErrCanceled) {
			_ = zenity.Error("Gagal memilih folder", zenity.Title("Error"))
		}
		return
	}

	files, err := utility.ValidatePDFFolder(folderPath)
	if err != nil {
		_ = zenity.Error(err.Error(), zenity.Title("Error"))
		return
	}

	if showConfirmationDialog(len(files)) {
		app.startProcessing(files)
	}
}

func (app *AppUI) selectFiles() {
	filePaths, err := zenity.SelectFileMultiple(
		zenity.Title("Pilih File PDF (Multiple)"),
		zenity.FileFilter{
			Name:     "PDF Files",
			Patterns: []string{"*.pdf"},
		},
	)
	if err != nil {
		if !errors.Is(err, zenity.ErrCanceled) {
			_ = zenity.Error("Gagal memilih file", zenity.Title("Error"))
		}
		return
	}

	if showConfirmationDialog(len(filePaths)) {
		app.startProcessing(filePaths)
	}
}

func (app *AppUI) startProcessing(files []string) {
	progress, err := zenity.Progress(
		zenity.Title("Memproses File PDF"),
		zenity.EntryText("Memulai pemrosesan..."),
		zenity.Modal(),
		zenity.Width(400),
		zenity.Height(150),
	)
	if err != nil {
		_ = zenity.Error("Gagal membuat progress dialog", zenity.Title("Error"))
		return
	}
	defer func(progress zenity.ProgressDialog) {
		err := progress.Close()
		if err != nil && !errors.Is(err, zenity.ErrCanceled) {
		}
	}(progress)

	app.logHndl.ResetLogs()

	ctx, cancel := context.WithCancel(context.Background())

	updateProgress := func(percent int, text string) error {
		if err := progress.Value(percent); err != nil {
			return err
		}
		return progress.Text(text)
	}

	go func() {
		<-progress.Done()
		cancel()
	}()

	result := core.ProcessFiles(ctx, app.logger, files, updateProgress)

	_ = progress.Value(100)
	_ = progress.Text("Pemrosesan selesai!")

	showResult(result)
}

func (app *AppUI) showLogs() {
	logs := app.logHndl.GetLogs()
	if len(logs) == 0 {
		_ = zenity.Info("Belum ada log yang tersedia.", zenity.Title("Informasi"))
		return
	}

	tempFile := filepath.Join(os.TempDir(), "sinta_log.txt")
	file, err := os.Create(tempFile)
	if err != nil {
		_ = zenity.Error("Gagal membuat file log sementara.", zenity.Title("Error"))
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
		}
	}()

	write := func(s string) {
		if err == nil {
			_, err = file.WriteString(s)
		}
	}

	write("=== SINTA ISSN CHECKER - LOG PEMROSESAN ===\n")
	write(fmt.Sprintf("Total Log Entries: %d\n\n", len(logs)))
	for i, message := range logs {
		write(fmt.Sprintf("[%03d] %s\n", i+1, message))
	}

	if err != nil {
		_ = zenity.Error("Gagal menulis ke file log sementara.", zenity.Title("Error"))
		return
	}

	err = zenity.Question(
		fmt.Sprintf("Log telah disimpan ke file:\n%s\n\nApakah Anda ingin membuka file log?", tempFile),
		zenity.Title("Log Pemrosesan"),
		zenity.OKLabel("Buka File"),
		zenity.CancelLabel("Tutup"),
	)

	if err == nil {
		openFile(tempFile)
	}
}

func showConfirmationDialog(fileCount int) bool {
	message := fmt.Sprintf("Terdapat %d file PDF.\n\nApakah Anda ingin melanjutkan proses?", fileCount)
	err := zenity.Question(
		message,
		zenity.Title("Konfirmasi"),
		zenity.OKLabel("Lanjutkan"),
		zenity.CancelLabel("Kembali"),
	)
	return err == nil
}

func showResult(result core.ProcessResult) {
	var b strings.Builder
	writeSection := func(title string, files []string) {
		if len(files) > 0 {
			_, _ = fmt.Fprintf(&b, "%s (%d file):\n", title, len(files))
			for _, file := range files {
				_, _ = fmt.Fprintf(&b, "   - %s\n", filepath.Base(file))
			}
			b.WriteString("\n")
		}
	}

	writeSection("BERHASIL DIUBAH", result.SuccessFiles)
	writeSection("GAGAL RENAME", result.FailedFiles)
	writeSection("ERROR PROSES", result.ProcessingErrorFiles)
	writeSection("TIDAK DITEMUKAN ISSN", result.ISSNNotFoundFiles)
	writeSection("TIDAK TERAKREDITASI", result.NotAccreditedFiles)
	writeSection("SUDAH SESUAI", result.AlreadyRenamedFiles)

	msg := b.String()
	if msg == "" {
		msg = "Tidak ada file yang diproses atau semua file gagal diproses."
	}

	isError := len(result.SuccessFiles) == 0 && (len(result.FailedFiles) > 0 || len(result.ProcessingErrorFiles) > 0)

	if isError {
		_ = zenity.Error(msg, zenity.Title("Pemrosesan Selesai dengan Masalah"), zenity.Width(700), zenity.Height(500))
	} else {
		_ = zenity.Info(msg, zenity.Title("Pemrosesan Selesai"), zenity.Width(700), zenity.Height(500))
	}
}

func openFile(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	_ = cmd.Run()
}
