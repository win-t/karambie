package middleware

import (
	"os"
	"path/filepath"

	"github.com/win-t/karambie"
	"github.com/win-t/karambie/middleware/logger"
	"github.com/win-t/karambie/middleware/recovery"
	"github.com/win-t/karambie/middleware/static"
)

func Common() karambie.HandlerList {
	detail := len(os.Getenv("PRODUCTION")) == 0
	curdir, _ := os.Getwd()
	return karambie.List(
		logger.Get(),
		recovery.Get(detail),
		static.Get(filepath.Join(curdir, "public")),
	)
}
