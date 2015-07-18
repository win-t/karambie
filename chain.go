package karambie

import (
	"net/http"
)

type Chain interface {
	Next() bool
	Stop()
	Resume() bool
}

func (c *context) prepare(h handlerList, r *http.Request) {
	c.hl = h
	c.req = r
	c.index = 0
	c.stoped = false
}

func (c *context) run() {
	for !c.stoped && c.index < len(c.hl) {
		c.hl[c.index].ServeHTTP(c, c.req)
		c.index += 1
		if c.written > 0 {
			c.stoped = true
		}
	}
}

func (c *context) Next() bool {
	if !c.stoped {
		c.index += 1
		c.run()
	}
	return !c.stoped
}

func (c *context) Stop() {
	c.stoped = true
}

func (c *context) Resume() bool {
	if c.stoped {
		c.stoped = false
		c.run()
	}
	return !c.stoped
}
