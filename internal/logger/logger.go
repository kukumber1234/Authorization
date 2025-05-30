package logger

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	file, err := os.OpenFile("C:\\Users\\Lenovo\\OneDrive\\Рабочий стол\\GTS\\Authorization\\internal\\logger\\log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed open log.txt", "error", err)
	}
	
	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}