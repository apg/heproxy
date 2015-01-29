package main

import (
	"log"
	"net/http"

	"github.com/apg/heproxy"
)

func main() {
	proxy := heproxy.NewHEProxy(nil)
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
