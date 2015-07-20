package main

import (
	"fmt"
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

func hello(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Hello World")
}

func sayhi(render render.Render, request *http.Request) {
	name := mux.Vars(request)["name"]
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
		// dont write anything, if you do that the chain list will be stop
	}
}

func main() {
	list := karambie.List()
	route := mux.NewRouter()       // use gorilla as router
	martini := martinihelper.New() // and also martini as middleware, see below

	currentDir, _ := os.Getwd()

	logger, log := logger.New(os.Stdout, "karambie")               // log every request
	recovery := recovery.New(true, log)                            // recover if panic
	notfoundhandler := notfoundhandler.New(true, nil)              // show 404, add trailing slash to url if necessary
	static := static.New(filepath.Join(currentDir, "public"), log) // serve static file in folder "public"

	// register logger service for martini
	martini.Map(log)

	//the list is immutable, every calling to Add will create new list
	// karambie.Later will create new handler that will be executed after succeeding handler
	list = list.Add(logger, recovery)                // list is [logger, recovery, ...]
	list = list.Add(karambie.Later(notfoundhandler)) // list is [logger, recovery, ..., notfoundhandler]
	list = list.Add(karambie.Later(static))          // list is [logger, recovery, ..., static, notfoundhandler]
	// or you can use karambie/middleware.Common() to build those list

	// list processing will stop if one of them respond the request
	// so, we can user list as NotFoundHandler (handle static file, and show error 404 if necessary)
	route.NotFoundHandler = list

	// see, the they are different list
	secureList := list.Add(martini.Conv(auth.Basic("user", "pass"))) // secureList is [logger, recovery, auth, ..., static, notfoundhandler]
	list = list.Add(martini.Conv(render.Renderer()))                 // list is       [logger, recovery, render, ..., static, notfoundhandler]

	// using http.HandlerFunc style
	// [logger, recovery, render, hello, (x)static, (x)notfoundhandler]
	// but 'static' and 'notfoundhandler' will be ignored because 'hello' response the request
	route.Handle("/helloworld", list.AddFunc(hello))

	// using martini.Handler style,
	// and gorilla routing
	// [logger, recovery, render, sayhi, (x)static, (x)notfoundhandler]
	route.Handle("/sayhi/{name}", list.Add(martini.Conv(sayhi)))

	// use secureList for sensitive resource
	// [logger, recovery, auth, secret, (x)static,(x)notfoundhandler]
	route.Handle("/secret", secureList.AddFunc(secret))

	// add filterheader to list
	apiList := list.AddFunc(filterheader)
	route.Handle("/api", apiList.AddFunc(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// get apikey
		ctx := karambie.Context(rw)
		key := ctx.Get("apikey").(string)
		fmt.Fprintln(rw, "Your api key : "+key)
	})))

	http.ListenAndServe(":3000", route)
}
