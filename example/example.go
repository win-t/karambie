package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/win-t/karambie"
	"github.com/win-t/karambie/middleware/logger"
	"github.com/win-t/karambie/middleware/notfoundhandler"
	"github.com/win-t/karambie/middleware/recovery"
	"github.com/win-t/karambie/middleware/static"
	"github.com/win-t/karambie/tool/martinihelper"

	"github.com/gorilla/mux"

	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"
)

func hello(rw http.ResponseWriter, r *http.Request) { // http.HandlerFunc style
	fmt.Fprintf(rw, "Hello World")
}

func sayhi(render render.Render, request *http.Request) { // martini.Handler style
	// get parameter from gorilla routing
	name := mux.Vars(request)["name"]

	// use martini-contrib/render
	render.JSON(http.StatusOK, "Hi "+name)
}

func secret(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "This is secret message")
}

func filterheader(rw http.ResponseWriter, r *http.Request) {
	if key := r.Header.Get("X-Api-Key"); len(key) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(rw, "X-Api-Key header not defined")
	} else {
		ctx := karambie.Context(rw)
		ctx.Set("apikey", key)
		// dont write anything, if you do that, the chain list will be stop
	}
}

func main() {
	list := karambie.List()
	route := mux.NewRouter()       // use gorilla as router
	martini := martinihelper.New() // and also martini as middleware, see below

	currentDir, _ := os.Getwd()

	log := log.New(os.Stdout, "[Karambie] ", 0)

	logger := logger.New(log)                                      // log every request
	recovery := recovery.New(nil, log)                             // recover if panic
	notfoundhandler := notfoundhandler.New(true, nil)              // show 404, add trailing slash to url if necessary
	static := static.New(filepath.Join(currentDir, "public"), log) // serve static file in folder "public"

	// register logger service for martini
	martini.Map(log)

	// the list is immutable, every calling to Add or AddFunc will create new list
	list = list.Add(logger, recovery)
	list = list.Add(karambie.Later(notfoundhandler))
	list = list.Add(karambie.Later(static))
	// or you can use karambie/middleware.Common() to build those list

	// list is [logger, recovery, notfoundhandler, static]
	// but the order of execution is [logger, recovery, static, notfoundhandler]
	// karambie.Later will create new handler that will be executed after succeeding handler
	// list processing will stop if one of them respond the request (http response status != 0)

	secureList := list.Add(martini.Conv(auth.Basic("user", "pass"))) // execution of secureList is [logger, recovery, auth, static, notfoundhandler]
	list = list.Add(martini.Conv(render.Renderer()))                 // execution of list is       [logger, recovery, render, static, notfoundhandler]
	// list != secureList, because it is immutable, every calling to Add or AddFunc will create new list

	// using http.HandlerFunc style
	route.Handle("/helloworld", list.AddFunc(hello)) // [logger, recovery, render, hello]
	// 'static' and 'notfoundhandler' will be ignored because 'hello' response the request

	// we can user list as NotFoundHandler (handle static file, and show error 404 if necessary)
	route.NotFoundHandler = list // [logger, recovery, static, notfoundhandler]
	// 'notfoundhandler' will be ignored if 'static' response the request

	// using martini.Handler style and gorilla routing
	route.Handle("/sayhi/{name}", list.Add(martini.Conv(sayhi))) // [logger, recovery, render, sayhi]

	// use secureList for sensitive resource
	route.Handle("/secret", secureList.AddFunc(secret)) // [logger, recovery, auth, secret]

	// add filterheader to list
	apiList := list.AddFunc(filterheader) // execution of apiList is [logger, recovery, filterheader, static, notfoundhandler]
	route.Handle("/api", apiList.AddFunc(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// this handler will not be called if 'filterheader' is not passed

		// get apikey
		ctx := karambie.Context(rw)
		key := ctx.Get("apikey").(string)

		fmt.Fprintln(rw, "Your api key : "+key)
	})))

	http.ListenAndServe(":3000", route)
}
