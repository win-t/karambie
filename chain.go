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
		c.hl[c.index].ServeHTTP(c, c.req)
		c.index += 1
		if c.written > 0 {
			c.stoped = true
		}
	}
}

func (c *ResponseWriterContext) Next() bool {
	if !c.stoped {
		c.index += 1
		c.run()
	}
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
