package notfoundhandler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/win-t/karambie"
)

const (
	template = `<html>
<head><title>Not Found</title>
<style type="text/css">
html, body {
	font-family: "Roboto", sans-serif;
	color: #333333;
	background-color: #ea5343;
	margin: 0px;
}
h1 {
	color: #d04526;
	background-color: #ffffff;
	padding: 20px;
	border-bottom: 1px dashed #2b3848;
}
pre {
	margin: 20px;
	padding: 20px;
	border: 2px solid #2b3848;
	background-color: #ffffff;
}
</style>
</head><body>
<h1>Not Found</h1>
<pre style="font-weight: bold;">%s</pre>
</body>
</html>`
)

func New(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c := karambie.Context(rw)
		if c.Status() == 0 {
			if strings.HasSuffix(r.URL.Path, "/") {
				c.WriteHeader(http.StatusNotFound)
				if h == nil {
					c.Write([]byte(fmt.Sprintf(template, r.URL)))
				} else {
					h.ServeHTTP(c, r)
				}
			} else {
				u, _ := url.Parse(r.URL.String())
				u.Path += "/"
				http.Redirect(rw, r, u.String(), http.StatusTemporaryRedirect)
			}
		}
	})
}
