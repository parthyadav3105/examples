package main

import (
	"log"
	"log/slog"

	"github.com/parthyadav3105/examples/pkg/logger"
)

func main() {

	err := logger.InitSlog(slog.LevelDebug, false)
	if err != nil {
		log.Fatalln("failed to setup slog with zap")
		return
	}
	defer logger.Cleanup() // for cleanup at zapLogger.Sync()

	// Normal slog usage:
	slog.Info("created vm at provider", slog.String("name", "demo"), slog.Int("vcpu", 8))
	slog.Warn("A warning occurred", slog.Any("error", "example error"))
	slog.Debug("read userinfo from db", slog.String("username", "ram"), slog.Any("password", logger.Secret("admin")))
	slog.Error("failed to attach volume at vm", slog.Any("error", err), slog.String("vm", "demo"))
}
