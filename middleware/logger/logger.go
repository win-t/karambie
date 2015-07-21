package logger

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/win-t/karambie"
)

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
// return http.Handler and log.Logger instance
func New(writer io.Writer, tag string) (http.Handler, *log.Logger) {
	log := log.New(writer, "["+tag+"] ", 0)
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

		c.Next()

		log.Printf(
			"%s %s for %s -> %v %s (written %d bytes) in %v\n",
			req.Method, req.URL.Path, addr,
			c.Status(), http.StatusText(c.Status()), c.Written(), time.Since(start),
		)
	}), log
}
