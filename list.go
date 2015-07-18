package karambie

import (
	"net/http"
)

type HandlerList []http.Handler

func List(h ...http.Handler) HandlerList {
	var ret HandlerList = make([]http.Handler, 0)
	return ret.Add(h...)
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

func (l HandlerList) AddFunc(h http.HandlerFunc) HandlerList {
	return l.Add(h)
}

func (h HandlerList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := Context(rw)
	c.prepare(h, r)
	c.run()
}
