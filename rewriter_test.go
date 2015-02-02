package heproxy

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"testing"
)

func TestLinearRewriter(t *testing.T) {
	EFFBytes := bytes.NewBufferString(EFF)
	ruleSet, _ := NewRuleSet(EFFBytes)
	rewriter := NewLinearRewriter([]*RuleSet{ruleSet})
	oldURL, _ := url.Parse("http://eff.org/welcome?funky=town")
	newURL := rewriter.Rewrite(oldURL)

	if newURL == oldURL || newURL.String() == oldURL.String() {
		t.Errorf("Exepected (%q) to be rewritten. Got: %q", oldURL.String(), newURL.String())
	}
}

func BenchmarkLinearRewriter(b *testing.B) {
	ruleSets := ReadRuleSets("./cmd/heproxy/rulesets")
	rewriter := NewLinearRewriter(ruleSets)

	rewriterBenchmark(ruleSets, rewriter, b)
}

func TestTrieRewriter(t *testing.T) {
	EFFBytes := bytes.NewBufferString(EFF)
	ruleSet, _ := NewRuleSet(EFFBytes)
	rewriter := NewTrieRewriter([]*RuleSet{ruleSet})
	oldURL, _ := url.Parse("http://foo.eff.org/welcome?funky=town")
	newURL := rewriter.Rewrite(oldURL)

	if newURL == oldURL || newURL.String() == oldURL.String() {
		t.Errorf("Exepected (%q) to be rewritten. Got: %q", oldURL.String(), newURL.String())
	}
}

func BenchmarkTrieRewriter(b *testing.B) {
	ruleSets := ReadRuleSets("./cmd/heproxy/rulesets")
	rewriter := NewTrieRewriter(ruleSets)

	rewriterBenchmark(ruleSets, rewriter, b)
}

func rewriterBenchmark(ruleSets []*RuleSet, rewriter Rewriter, b *testing.B) {
	// Build 1000 URLs at random to attempt to rewrite.
	urls := buildTestUrls(ruleSets, 1000)

	for i := 0; i < b.N; i++ {
		for _, u := range urls {
			rewriter.Rewrite(u)
		}
	}
}

func buildTestUrls(rss []*RuleSet, cnt int) []*url.URL {
	urls := make([]*url.URL, cnt)
	for i := 0; i < cnt; i++ {
		// select a random ruleset
	again:
		n := rand.Int31n(int32(cnt))
		rs := rss[n]
		if len(rs.Targets) == 0 {
			goto again
		}

		urls[i], _ = url.Parse(fmt.Sprintf("http://%s/test/url",
			strings.Replace(rs.Targets[0].Host, "*", "www", -1)))
	}

	return urls
}
