package karambie

import (
	"net/http"
)

const maxInt int = int(^uint(0) >> 1)

type ResponseWriterContext struct {
	rw      http.ResponseWriter
	status  int
	written int

	data map[interface{}]interface{}

	hl     HandlerList
	req    *http.Request
	index  int
	stoped bool
}

func Context(rw http.ResponseWriter) (ret *ResponseWriterContext) {
	ret, ok := rw.(*ResponseWriterContext)
	if !ok {
		ret = &ResponseWriterContext{rw, 0, 0, make(map[interface{}]interface{}), nil, nil, maxInt, true}
	}
	return
}

func (c *ResponseWriterContext) Header() http.Header {
	return c.rw.Header()
}

func (c *ResponseWriterContext) Write(b []byte) (int, error) {
	if c.status == 0 {
		c.WriteHeader(http.StatusOK)
	}
	size, err := c.rw.Write(b)
	c.written += size
	return size, err
}

func (c *ResponseWriterContext) WriteHeader(s int) {
	if c.status != 0 {
		return
	}
	c.rw.WriteHeader(s)
	c.status = s
}

func (c *ResponseWriterContext) Original() http.ResponseWriter {
	return c.rw
}

func (c *ResponseWriterContext) Status() int {
	return c.status
}

func (c *ResponseWriterContext) Written() int {
	return c.written
}

func (c *ResponseWriterContext) Set(key, value interface{}) {
	c.data[key] = value
}

func (c *ResponseWriterContext) GetOk(key interface{}) (ret interface{}, ok bool) {
	ret, ok = c.data[key]
	return
}

func (c *ResponseWriterContext) Get(key interface{}) (ret interface{}) {
	ret, _ = c.GetOk(key)
	return
}

func (c *ResponseWriterContext) Delete(key interface{}) {
	delete(c.data, key)
}
