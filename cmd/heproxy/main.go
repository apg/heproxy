package main

import (
	"log"
	"net/http"
	"os"

	"github.com/apg/heproxy"
)

func readRuleSets() []*heproxy.RuleSet {
	var ruleSets []*heproxy.RuleSet

	dir, err := os.Open("./rulesets")
	if err != nil {
		panic(err)
	}

	dirnames, err := dir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	for _, fn := range dirnames {
		f, _ := os.Open("./rulesets/" + fn)
		rs, err := heproxy.NewRuleSet(f)
		if err != nil {
			log.Printf("err=%q ruleset=%q", err, f)
			f.Close()
			continue
		} else {
			ruleSets = append(ruleSets, rs)
		}
		f.Close()
	}

	return ruleSets
}

func main() {
	rewriter := heproxy.NewRewriter(readRuleSets())

	proxy := heproxy.NewHEProxy(rewriter)
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
