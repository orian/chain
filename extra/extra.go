package chain

import (
	"github.com/orian/wctx"
	"golang.org/x/net/context"

	"net/http"
	"strings"
)

// StripPrefix returns a handler that serves HTTP requests
// by removing the given prefix from the request URL's Path
// and invoking the handler h. StripPrefix handles a
// request for a path that doesn't begin with prefix by
// replying with an HTTP 404 not found error.
func StripPrefix(prefix string) func(wctx.Handler) wctx.Handler {
	return func(h wctx.Handler) wctx.Handler {
		if prefix == "" {
			return h
		}
		return wctx.HandleFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
				r.URL.Path = p
				h.ServeHTTP(c, w, r)
			}
			// else {
			// 	// problematic case
			// 	NotFound(w, r)
			// }
		})
	}
}
