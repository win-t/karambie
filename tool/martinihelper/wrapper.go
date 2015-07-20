package martinihelper

import (
	"net/http"
	"reflect"

	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/win-t/karambie"
)

type key int

const (
	contextInstance key = iota
)

type MartiniHelper struct {
	inject.Injector
}

func New() *MartiniHelper {
	this := &MartiniHelper{inject.New()}
	retHandler := martini.New().Get(reflect.TypeOf(martini.ReturnHandler(nil))).Interface()
	// retHandler := martini.defaultReturnHandler()
	this.Map(retHandler)
	return this
}

func (this *MartiniHelper) Conv(h martini.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rwc := karambie.Context(rw)
		c := this.context(rwc, r)

		vals, err := c.Invoke(h)
		if err != nil {
			panic(err)
		}

		if rwc.Status() == 0 {
			// if the handler returned something, write it to the http response
			if len(vals) > 0 {
				ev := c.Get(reflect.TypeOf(martini.ReturnHandler(nil)))
				handleReturn := ev.Interface().(martini.ReturnHandler)
				handleReturn(c, vals)
			}
		}
	})
}

type context struct {
	inject.Injector
}

func (this *MartiniHelper) context(rwc *karambie.ResponseWriterContext, r *http.Request) martini.Context {
	if v, ok := rwc.GetOk(contextInstance); ok {
		return v.(martini.Context)
	} else {
		c := &context{inject.New()}
		c.SetParent(this)

		c.Map(rwc)
		c.MapTo(c, (*martini.Context)(nil))
		c.MapTo(rwc, (*http.ResponseWriter)(nil))
		c.Map(r)

		rwc.Set(contextInstance, c)
		return c
	}
}

func (c *context) rwc() *karambie.ResponseWriterContext {
	return c.Get(reflect.TypeOf((*karambie.ResponseWriterContext)(nil))).Interface().(*karambie.ResponseWriterContext)
}

func (c *context) Next() {
	c.rwc().Next()
}

func (c *context) Written() bool {
	return c.rwc().Status() != 0
}
