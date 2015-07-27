package karambie

import (
	"net/http"
	"sync"
)

// Context is created and destroyed for each request/response
type ResponseWriterContext struct {
	rw http.ResponseWriter

	status  int
	written int

	data map[interface{}]interface{}
	sync sync.RWMutex

	hl      HandlerList
	req     *http.Request
	index   int
	running bool
}

// create new context of http.ResponseWriter, or just return the input if it's already have context
func Context(rw http.ResponseWriter) (ret *ResponseWriterContext) {
	ret, ok := rw.(*ResponseWriterContext)
	if !ok {
		ret = &ResponseWriterContext{
			rw:   rw,
			data: make(map[interface{}]interface{}),
		}
	}
	return
}

// see http.ResponseWriter
func (c *ResponseWriterContext) Header() http.Header {
	return c.rw.Header()
}

// see http.ResponseWriter
func (c *ResponseWriterContext) Write(b []byte) (int, error) {
	if c.status == 0 {
		// c.rw.Write(b) will set status to http.StatusOK
		c.status = http.StatusOK
	}
	size, err := c.rw.Write(b)
	c.written += size
	return size, err
}

// see http.ResponseWriter
func (c *ResponseWriterContext) WriteHeader(s int) {
	if c.status != 0 {
		return
	}
	c.rw.WriteHeader(s)
	c.status = s
}

// return original http.ResponseWriter (without context)
func (c *ResponseWriterContext) Original() http.ResponseWriter {
	return c.rw
}

// HTTP Response status, 0 means no response at all
func (c *ResponseWriterContext) Status() int {
	return c.status
}

// size in byte that has been written
func (c *ResponseWriterContext) Written() int {
	return c.written
}

// set data in context, identified with key
func (c *ResponseWriterContext) Set(key, value interface{}) {
	c.sync.Lock()
	defer c.sync.Unlock()

	c.data[key] = value
}

// get data in context, identified with key, will return (nil, false) if data doesn't exist
func (c *ResponseWriterContext) GetOk(key interface{}) (ret interface{}, ok bool) {
	c.sync.RLock()
	defer c.sync.RUnlock()

	ret, ok = c.data[key]
	return
}

// get data in context, identified with key, will return nil if data doesn't exist
func (c *ResponseWriterContext) Get(key interface{}) (ret interface{}) {
	ret, _ = c.GetOk(key)
	return
}

// delete data in context, identified with key
func (c *ResponseWriterContext) Delete(key interface{}) {
	c.sync.Lock()
	defer c.sync.Unlock()

	delete(c.data, key)
}
