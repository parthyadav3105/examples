// Logger package to initialize slog with zap.
//
// Log levels in decreasing order of precedence:
// Error > Warn > Info > Debug
//
// Example Usage:
//
//	 func main() {
//
//		err := logger.InitSlog(slog.LevelDebug, false)
//		if err != nil {
//			log.Fatalln("failed to setup slog with zap")
//			return
//		}
//		defer logger.Cleanup() // for cleanup at zapLogger.Sync()
//
//			// Normal slog usage:
//			slog.Info("created vm at provider", slog.String("name", "demo"), slog.Int("vcpu", 8))
//			slog.Debug("read userinfo from db", slog.String("username", "ram"), slog.Any("password", logger.Secret("admin")))
//			slog.Error("failed to attach volume at vm", slog.Any("error", err), slog.String("vm", "demo"))
//		}
//
// Output:
//
//	{"level":"info","time":"2025-11-08T21:31:09.311+0530","line":"inf/main.go:19","msg":"created vm at provider","vcpu":8,"name":"demo"}
//	{"level":"debug","time":"2025-11-08T21:31:09.311+0530","line":"inf/main.go:20","msg":"read userinfo from db","username":"ram","password":"****"}
//	{"level":"error","time":"2025-11-08T21:31:09.311+0530","line":"inf/main.go:21","msg":"failed to attach volume at vm","vm":"demo"}
package logger

import (
	"log/slog"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Use `logger.Secret("mycred")` to redact secret values when logged.
type Secret string

func (s Secret) LogValue() slog.Value {
	return slog.StringValue("****")
}

// Since slogzap (Zap adapter for slog) does not honor LogValue() above, we use
// Stringer to hide secrets. The LogValue method above is still implemented as a
// fallback for other slog handlers that do respect it.
func (s Secret) String() string {
	return "****"
}

var zapLogger *zap.Logger

// Init Slog with Zap for production
func InitSlog(level slog.Leveler, developmentMode bool) error {
	// Create a zap logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableCaller = false
	config.EncoderConfig.CallerKey = "line"

	// fix: slogzap.Option.Level does not seems to respect slog.LevelDebug, below line fixes it.
	config.Level = zap.NewAtomicLevelAt(slogzap.LogLevels[level.Level()])

	if developmentMode {
		config.Development = true
		config.Encoding = "console"
		// config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	zapLogger, err := config.Build()
	if err != nil {
		return err
	}

	// Create a slog handler using the zap logger
	handler := slogzap.Option{
		Logger:    zapLogger,
		Level:     level, // Set the desired slog level
		AddSource: true,
	}.NewZapHandler()

	// Set as default slog logger
	slog.SetDefault(slog.New(handler))
	return nil
}

func Cleanup() {
	if zapLogger != nil {
		zapLogger.Sync() // Flushes any buffered log entries
	}
}
