package karambie

import (
	"net/http"
)

type HandlerList []http.Handler

func convFuncToHandler(h []http.HandlerFunc) []http.Handler {
	hs := make([]http.Handler, 0, len(h))
	for _, v := range h {
		hs = append(hs, v)
	}
	return hs
}

func List(h ...http.Handler) HandlerList {
	var ret HandlerList = make([]http.Handler, 0)
	return ret.Add(h...)
}

func ListFunc(h ...http.HandlerFunc) HandlerList {
	return List(convFuncToHandler(h)...)
}

func (l HandlerList) Add(h ...http.Handler) HandlerList {
	var ret HandlerList = append(make([]http.Handler, 0, len(l)+len(h)), l...)
	for _, v := range h {
		if hl, ok := v.(HandlerList); ok {
			ret = append(ret, hl...)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}

func (l HandlerList) AddFunc(h ...http.HandlerFunc) HandlerList {
	return l.Add(convFuncToHandler(h)...)
}

func (h HandlerList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := Context(rw)
	c.prepare(h, r)
	c.run()
}

func Pending(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c := Context(rw)
		if c.Next() {
			h.ServeHTTP(rw, r)
		}
	})
}

func PendingFunc(h http.HandlerFunc) http.HandlerFunc {
	return Pending(h).ServeHTTP
}
