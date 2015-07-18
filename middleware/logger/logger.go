package logger

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/win-t/karambie"
)

type key int

const instance key = 0

func Current(c karambie.ResponseWriterContext) *log.Logger {
	if v, ok := c.GetOk(instance); ok {
		if v, ok := v.(*log.Logger); ok {
			return v
		}
	}
	return nil
}

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Get() http.Handler {
	log := log.New(os.Stdout, "[Karambie] ", 0)
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		c := karambie.Context(res)
		c.Set(instance, log)

		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		log.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

		c.Next()

		log.Printf("Completed %v %s in %v\n", c.Status(), http.StatusText(c.Status()), time.Since(start))
	})
}
