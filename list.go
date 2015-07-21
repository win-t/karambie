package karambie

import (
	"net/http"
)

// group multiple http.Handler into one http.Handler
type HandlerList []http.Handler

// convert []http.HandlerFunc into []http.Handler
func ConvList(h []http.HandlerFunc) []http.Handler {
	hs := make([]http.Handler, 0, len(h))
	for _, v := range h {
		hs = append(hs, v)
	}
	return hs
}

// create new HandlerList from http.Handler
func List(h ...http.Handler) HandlerList {
	var ret HandlerList = make([]http.Handler, 0)
	return ret.Add(h...)
}

// create new HandlerList from http.HandlerFunc
func ListFunc(h ...http.HandlerFunc) HandlerList {
	return List(ConvList(h)...)
}

// create new HandlerList and add http.Handler to it
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

// create new HandlerList and add http.HandlerFunc to it
func (l HandlerList) AddFunc(h ...http.HandlerFunc) HandlerList {
	return l.Add(ConvList(h)...)
}

// see http.Hadler
func (h HandlerList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := Context(rw)
	c.prepare(h, r)
	c.run()
}

// create new http.Handler wrapper that will be executed later in list
func Later(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c := Context(rw)
		if c.Next() {
			// run h if list chain still active, i.e. no one response the request
			h.ServeHTTP(rw, r)
		}
	})
}

// create new http.HandlerFunc wrapper that will be executed later in list
func LaterFunc(h http.HandlerFunc) http.HandlerFunc {
	return Later(h).ServeHTTP
}
