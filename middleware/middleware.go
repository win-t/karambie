package middleware

import (
	"log"
	"os"
	"path/filepath"

	"github.com/win-t/karambie"
	"github.com/win-t/karambie/middleware/logger"
	"github.com/win-t/karambie/middleware/notfoundhandler"
	"github.com/win-t/karambie/middleware/recovery"
	"github.com/win-t/karambie/middleware/static"
)

// get common HandlerList, it contain [logger, recovery, notfoundhandler, static]
func Common() (karambie.HandlerList, *log.Logger) {
	tag := filepath.Base(os.Args[0])
	log := log.New(os.Stdout, "["+tag+"] ", 0)
	cwd, _ := os.Getwd()

	list := karambie.List(
		logger.New(log),
		recovery.New(log, nil),
		karambie.Later(notfoundhandler.New(true, nil)),
		karambie.Later(static.New(filepath.Join(cwd, "public"), log)),
	)
	return list, log
}
