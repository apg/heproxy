package heproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func httpsDirector(req *http.Request, rules *RuleSet) {
	log.Printf("%s", req.URL)
}

// NewHEProxy returns a new proxy which maps a request to an HTTPS url if available
func NewHEProxy(rs *RuleSet) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		httpsDirector(req, rs)
	}

	return &httputil.ReverseProxy{Director: director}
}
