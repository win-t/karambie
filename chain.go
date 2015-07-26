package karambie

import (
	"net/http"
)

// prepare context with HandlerList
func (c *ResponseWriterContext) prepare(h HandlerList, r *http.Request) {
	c.hl = h
	c.req = r
	c.index = 0
	c.running = true
}

// run context
func (c *ResponseWriterContext) run() {
	for c.running && c.index < len(c.hl) {
		h := c.hl[c.index]
		c.index += 1

		h.ServeHTTP(c, c.req)

		// status != 0, mean h handle the request
		if c.status != 0 {
			c.Stop()
		}
	}
}

// calling to this method will block until succeeding handler in the list is executed.
// return false if list execution is stopped
func (c *ResponseWriterContext) Next() bool {
	c.run()
	return c.running
}

// Stop the list execution. it will be stopped if one of handler in the list response the request (set http response status to something that not 0)
// e.g. calling ResponseWritter.WriteHeader or ResponseWritter.Write
func (c *ResponseWriterContext) Stop() {
	c.running = false
}

// Resume list execution, use this method after calling 'Next' method.
// Calling to this method without calling 'Next' is no-op
func (c *ResponseWriterContext) Resume() bool {
	if !c.running {
		c.running = true
		c.run()
	}
	return c.running
}
