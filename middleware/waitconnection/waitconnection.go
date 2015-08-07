package waitconnection

import (
	"net/http"
	"sync"

	"github.com/win-t/karambie"
)

func New() (http.Handler, *sync.WaitGroup) {
	var waiter sync.WaitGroup

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		waiter.Add(1)
		defer waiter.Done()
		karambie.Context(rw).Next()
	}), &waiter
}
