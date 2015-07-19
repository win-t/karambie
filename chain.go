package karambie

import (
	"net/http"
)

func (c *ResponseWriterContext) prepare(h HandlerList, r *http.Request) {
	c.hl = h
	c.req = r
	c.index = 0
	c.stoped = false
}

func (c *ResponseWriterContext) run() {
	for !c.stoped && c.index < len(c.hl) {
		h := c.hl[c.index]
		c.index += 1

		h.ServeHTTP(c, c.req)
		if c.status != 0 {
			c.stoped = true
		}
	}
}

func (c *ResponseWriterContext) Next() bool {
	c.run()
	return !c.stoped
}

func (c *ResponseWriterContext) Stop() {
	c.stoped = true
}

func (c *ResponseWriterContext) Resume() bool {
	if c.stoped {
		c.stoped = false
		c.run()
	}
	return !c.stoped
}
