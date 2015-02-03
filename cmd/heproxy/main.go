package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/apg/heproxy"
)

var (
	listen = flag.String("listen", ":9999", "address:port to listen on")
	help   = flag.Bool("help", false, "this help message")
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	ruleSets := heproxy.ReadRuleSets(path.Join(
		os.Getenv("GOPATH"), "src", "github.com", "apg", "heproxy", "cmd", "heproxy", "rulesets"))
	rewriter := heproxy.NewLinearRewriter(ruleSets)
	proxy := heproxy.NewHEProxy(rewriter)
	log.Fatal(http.ListenAndServe(*listen, proxy))
}
