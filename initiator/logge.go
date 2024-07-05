package initiator

import (
	"os"

	"golang.org/x/exp/slog"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
