package utility

import (
	"context"
	"log/slog"
	"strings"
	"sync"
)

type CaptureHandler struct {
	handler     slog.Handler
	logMessages []string
	logMutex    sync.Mutex
}

func NewCaptureHandler(handler slog.Handler) *CaptureHandler {
	return &CaptureHandler{
		handler:     handler,
		logMessages: []string{},
	}
}

func (h *CaptureHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.handler == nil {
		return false
	}
	return h.handler.Enabled(ctx, level)
}

func (h *CaptureHandler) Handle(ctx context.Context, r slog.Record) error {
	var buf strings.Builder
	buf.WriteString(r.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString(" [")
	buf.WriteString(r.Level.String())
	buf.WriteString("] ")
	buf.WriteString(r.Message)

	if r.NumAttrs() > 0 {
		buf.WriteString(" ")
		r.Attrs(func(a slog.Attr) bool {
			buf.WriteString(a.Key)
			buf.WriteString("=")
			buf.WriteString(a.Value.String())
			return true
		})
	}

	h.logMutex.Lock()
	h.logMessages = append(h.logMessages, buf.String())
	h.logMutex.Unlock()

	if h.handler != nil {
		return h.handler.Handle(ctx, r)
	}
	return nil
}

func (h *CaptureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.handler == nil {
		return NewCaptureHandler(nil)
	}
	return NewCaptureHandler(h.handler.WithAttrs(attrs))
}

func (h *CaptureHandler) WithGroup(name string) slog.Handler {
	if h.handler == nil {
		return NewCaptureHandler(nil)
	}
	return NewCaptureHandler(h.handler.WithGroup(name))
}

func (h *CaptureHandler) GetLogs() []string {
	h.logMutex.Lock()
	defer h.logMutex.Unlock()
	logsCopy := make([]string, len(h.logMessages))
	copy(logsCopy, h.logMessages)
	return logsCopy
}

func (h *CaptureHandler) ResetLogs() {
	h.logMutex.Lock()
	h.logMessages = []string{}
	h.logMutex.Unlock()
}
