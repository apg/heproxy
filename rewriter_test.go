package heproxy

import (
	"bytes"
	"log"
	"net/url"
	"testing"
)

func TestRewriter(t *testing.T) {
	EFFBytes := bytes.NewBufferString(EFF)
	ruleSet, _ := NewRuleSet(EFFBytes)
	rewriter := NewRewriter([]*RuleSet{ruleSet})
	oldURL, _ := url.Parse("http://eff.org/welcome?funky=town")
	newURL := rewriter.Rewrite(oldURL)

	if newURL == oldURL || newURL.String() == oldURL.String() {
		t.Errorf("Exepected (%q) to be rewritten. Got: %q", oldURL.String(), newURL.String())
	}

	log.Println(newURL)
}
