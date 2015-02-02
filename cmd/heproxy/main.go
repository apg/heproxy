package main

import (
	"log"
	"net/http"

	"github.com/apg/heproxy"
)

func main() {
	ruleSets := heproxy.ReadRuleSets("./rulesets")
	rewriter := heproxy.NewRewriter(ruleSets)
	proxy := heproxy.NewHEProxy(rewriter)
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
