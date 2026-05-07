package log

import (
	"log/slog"
	"os"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

func New(version string) *slog.Logger {
	return slogutil.New(&slogutil.InputNew{
		Name:    "tfprovidercheck",
		Version: version,
		Out:     os.Stderr,
	}).Logger
}
