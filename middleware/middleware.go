package middleware

import (
	"log"
	"os"

	"github.com/win-t/karambie"
	"github.com/win-t/karambie/middleware/logger"
	"github.com/win-t/karambie/middleware/notfoundhandler"
	"github.com/win-t/karambie/middleware/recovery"
	"github.com/win-t/karambie/middleware/static"
)

func Common(staticDir string) (karambie.HandlerList, *log.Logger) {
	logger, l := logger.New(os.Stdout, "Karambie")
	ret := karambie.List(
		logger,
		recovery.New(true, l),
		karambie.Pending(notfoundhandler.New(nil)),
	)
	if len(staticDir) > 0 {
		ret = ret.Add(karambie.Pending(static.New(staticDir, l)))
	}
	return ret, l
}
