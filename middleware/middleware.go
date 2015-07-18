package middleware

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/win-t/karambie/middleware/logger"
	"github.com/win-t/karambie/middleware/recovery"
	"github.com/win-t/karambie/middleware/static"
)

func Common() []http.Handler {
	detail := len(os.Getenv("PRODUCTION")) == 0
	curdir, _ := os.Getwd()
	return []http.Handler{
		logger.Get(),
		recovery.Get(detail),
		static.Get(filepath.Join(curdir, "public")),
	}
}
