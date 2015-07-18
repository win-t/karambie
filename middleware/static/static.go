package static

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/win-t/karambie"
	"github.com/win-t/karambie/middleware/logger"
)

// StaticOptions is a struct for specifying configuration options for the martini.Static middleware.
type StaticOptions struct {
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// SkipLogging will disable [Static] log messages when a static file is served.
	SkipLogging bool
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
	// Fallback defines a default URL to serve when the requested resource was
	// not found.
	Fallback string
	// Exclude defines a pattern for URLs this handler should never process.
	Exclude string
}

func prepareStaticOptions(options []StaticOptions) StaticOptions {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}
	// Normalize the prefix if provided
	if opt.Prefix != "" {
		// Ensure we have a leading '/'
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		// Remove any trailing '/'
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
	return opt
}

// Static returns a middleware handler that serves static files in the given directory.
func Get(directory string, staticOpt ...StaticOptions) http.Handler {
	if !filepath.IsAbs(directory) {
		panic(directory + " must be absolute path")
	}
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		c := karambie.Context(res)
		log := logger.Current(c)

		if req.Method != "GET" && req.Method != "HEAD" {
			return
		}
		if opt.Exclude != "" && strings.HasPrefix(req.URL.Path, opt.Exclude) {
			return
		}
		file := req.URL.Path
		// if we have a prefix, filter requests by stripping the prefix
		if opt.Prefix != "" {
			if !strings.HasPrefix(file, opt.Prefix) {
				return
			}
			file = file[len(opt.Prefix):]
			if file != "" && file[0] != '/' {
				return
			}
		}
		f, err := dir.Open(file)
		if err != nil {
			// try any fallback before giving up
			if opt.Fallback != "" {
				file = opt.Fallback // so that logging stays true
				f, err = dir.Open(opt.Fallback)
			}

			if err != nil {
				// discard the error?
				return
			}
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return
		}

		// try to serve index file
		if fi.IsDir() {
			// redirect if missing trailing slash
			if !strings.HasSuffix(req.URL.Path, "/") {
				dest := url.URL{
					Path:     req.URL.Path + "/",
					RawQuery: req.URL.RawQuery,
					Fragment: req.URL.Fragment,
				}
				http.Redirect(res, req, dest.String(), http.StatusFound)
				return
			}

			file = path.Join(file, opt.IndexFile)
			f, err = dir.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				return
			}
		}

		if !opt.SkipLogging && log != nil {
			log.Println("[Static] Serving " + file)
		}

		// Add an Expires header to the static content
		if opt.Expires != nil {
			res.Header().Set("Expires", opt.Expires())
		}

		http.ServeContent(res, req, file, fi.ModTime(), f)
	})
}
