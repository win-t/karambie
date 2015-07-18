package karambie

import (
	"net/http"
)

const maxInt int = int(^uint(0) >> 1)

type ResponseWriterContext interface {
	http.ResponseWriter
	Original() http.ResponseWriter

	Status() int
	Written() int

	Set(interface{}, interface{})
	GetOk(interface{}) (interface{}, bool)
	Get(interface{}) interface{}
	Delete(interface{})

	Chain
}

type context struct {
	rw      http.ResponseWriter
	status  int
	written int

	data map[interface{}]interface{}

	hl     handlerList
	req    *http.Request
	index  int
	stoped bool
}

func Context(rw http.ResponseWriter) (ret ResponseWriterContext) {
	ret, ok := rw.(ResponseWriterContext)
	if !ok {
		ret = &context{rw, 0, 0, make(map[interface{}]interface{}), nil, nil, maxInt, true}
	}
	return
}

func (c *context) Header() http.Header {
	return c.rw.Header()
}

func (c *context) Write(b []byte) (int, error) {
	if c.status == 0 {
		c.WriteHeader(http.StatusOK)
	}
	size, err := c.rw.Write(b)
	c.written += size
	return size, err
}

func (c *context) WriteHeader(s int) {
	if c.status != 0 {
		return
	}
	c.rw.WriteHeader(s)
	c.status = s
}

func (c *context) Original() http.ResponseWriter {
	return c.rw
}

func (c *context) Status() int {
	return c.status
}

func (c *context) Written() int {
	return c.written
}

func (c *context) Set(key, value interface{}) {
	c.data[key] = value
}

func (c *context) GetOk(key interface{}) (ret interface{}, ok bool) {
	ret, ok = c.data[key]
	return
}

func (c *context) Get(key interface{}) (ret interface{}) {
	ret, _ = c.GetOk(key)
	return
}

func (c *context) Delete(key interface{}) {
	delete(c.data, key)
}
