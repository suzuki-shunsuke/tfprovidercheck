package log

import (
	"log/slog"
	"os"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

func New(version string) (*slog.Logger, *slog.LevelVar) {
	logLevelVar := &slog.LevelVar{}
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "tfprovidercheck",
		Version: version,
		Out:     os.Stderr,
		Level:   logLevelVar,
	})
	return logger, logLevelVar
}

func SetLevel(logLevelVar *slog.LevelVar, level string) error {
	if level == "" {
		return nil
	}
	return slogutil.SetLevel(logLevelVar, level) //nolint:wrapcheck
}
