package logger

import (
	"log"
	"net/http"
	"time"

	"github.com/win-t/karambie"
)

func logPrintf(l *log.Logger, format string, v ...interface{}) {
	if l == nil {
		log.Printf(format, v...)
	} else {
		l.Printf(format, v...)
	}
}

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
// return http.Handler
func New(log *log.Logger, excludeHttpOk bool) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		c := karambie.Context(res)

		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		defer func() {
			if err := recover(); err != nil {
				logPrintf(log, "UNHANDLED PANIC")
				panic(err)
			}
		}()

		c.Next()

		if excludeHttpOk && c.Status() == http.StatusOK {
			return
		}

		logPrintf(log,
			"%s %s for %s -> %v %s (written %d bytes) in %v\n",
			req.Method, req.URL.Path, addr,
			c.Status(), http.StatusText(c.Status()), c.Written(), time.Since(start),
		)
	})
}
