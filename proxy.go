package heproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// NewHEProxy returns a new proxy which maps a request to an HTTPS url if available
func NewHEProxy(rw Rewriter) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		oldURL := req.URL
		newURL := rw.Rewrite(oldURL)
		req.URL = newURL
		log.Printf("action=rewrite from=%q to=%q", oldURL, newURL)
	}

	return &httputil.ReverseProxy{Director: director}
}
