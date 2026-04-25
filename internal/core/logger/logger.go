package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Log — глобальный логер. Используй везде: logger.Log.Info(...)
var Log *slog.Logger

// ── ANSI цвета ────────────────────────────────────────────────────────────────

const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorCyan // Debug
	}
}

func levelLabel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERR"
	case level >= slog.LevelWarn:
		return "WRN"
	case level >= slog.LevelInfo:
		return "INF"
	default:
		return "DBG"
	}
}

// ── Красивый handler для консоли ─────────────────────────────────────────────

type prettyHandler struct {
	out  io.Writer
	opts slog.HandlerOptions
}

func newPrettyHandler(out io.Writer, opts slog.HandlerOptions) *prettyHandler {
	return &prettyHandler{out: out, opts: opts}
}

func (h *prettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *prettyHandler) WithGroup(name string) slog.Handler       { return h }

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	color := levelColor(r.Level)
	label := levelLabel(r.Level)

	// Время — серое
	ts := colorGray + r.Time.Format("2006-01-02 15:04:05") + colorReset

	// Уровень — цветной, жирный, в скобках
	lvl := fmt.Sprintf("%s%s[%s]%s", colorBold, color, label, colorReset)

	// Сообщение — белое жирное
	msg := fmt.Sprintf("%s%s%s", colorBold, r.Message, colorReset)

	// Атрибуты — key=value серым цветом
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s%s%s=%v", colorGray, a.Key, colorReset, a.Value))
		return true
	})

	line := fmt.Sprintf("%s %s %s", ts, lvl, msg)
	if len(attrs) > 0 {
		line += "  " + strings.Join(attrs, "  ")
	}
	line += "\n"

	_, err := fmt.Fprint(h.out, line)
	return err
}

// ── Init ─────────────────────────────────────────────────────────────────────

// Init инициализирует логер. Вызови один раз в main().
// logsDir — папка для файлов, например "logs"
func Init(logsDir string) error {
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return err
	}

	fileName := filepath.Join(logsDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}

	// Консоль — красивый цветной текст
	consoleHandler := newPrettyHandler(os.Stdout, opts)

	// Файл — JSON (удобно для grep/парсинга)
	fileHandler := slog.NewJSONHandler(file, &opts)

	// Оба handler'а одновременно
	Log = slog.New(&multiHandler{
		handlers: []slog.Handler{consoleHandler, fileHandler},
	})

	slog.SetDefault(Log)
	return nil
}

// ── multiHandler — пишет в несколько handler'ов сразу ────────────────────────

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}