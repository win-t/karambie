package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/win-t/karambie"
	"github.com/win-t/karambie/middleware/logger"
	"github.com/win-t/karambie/middleware/notfoundhandler"
	"github.com/win-t/karambie/middleware/recovery"
	"github.com/win-t/karambie/middleware/static"
)

func Common(verboseError bool, notFound http.Handler, staticDir string) (karambie.HandlerList, *log.Logger) {
	logger, l := logger.New(os.Stdout, "Karambie")
	ret := karambie.List(
		logger,
		recovery.New(verboseError, l),
		karambie.Later(notfoundhandler.New(notFound)),
	)
	if len(staticDir) > 0 {
		ret = ret.Add(karambie.Later(static.New(staticDir, l)))
	}
	return ret, l
}
