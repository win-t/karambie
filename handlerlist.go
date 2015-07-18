package karambie

import (
	"net/http"
)

type handlerList []http.Handler

func BuildChain(h []http.Handler) handlerList {
	ret := handlerList(make([]http.Handler, 0))
	for _, v := range h {
		if hl, ok := v.(handlerList); ok {
			ret = append(ret, []http.Handler(hl)...)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func (h handlerList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := Context(rw).(*context)
	c.prepare(h, r)
	c.run()
}
