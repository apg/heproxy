package heproxy

import (
	"bytes"
	"net/url"
	"testing"
)

const EFF = `<ruleset name="EFF">
  <target host="eff.org" />
  <target host="*.eff.org" />

  <exclusion pattern="^http://action\.eff\.org/"/>
  <rule from="^http://eff\.org/" to="https://eff.org/"/>
  <rule from="^http://([^/:@]*)\.eff\.org/" to="https://$1.eff.org/"/>
</ruleset>`

func TestRuleSetReader(t *testing.T) {
	EFFBytes := bytes.NewBufferString(EFF)
	ruleSet, err := NewRuleSet(EFFBytes)
	if err != nil {
		t.Errorf("NewRuleSet failed: %q", err)
	}

	if !ruleSet.Exclusions[0].patternRE.MatcherString("http://action.eff.org/", 0).Matches() {
		t.Errorf("Expected http://action.eff.org to match against exclusion")
	}
}

func TestRuleApply(t *testing.T) {
	EFFBytes := bytes.NewBufferString(EFF)
	ruleSet, _ := NewRuleSet(EFFBytes)

	oldURL, _ := url.Parse("http://foo.eff.org/welcome?funky=town")
	newURL := ruleSet.Rewrite(oldURL)

	if newURL.String() == oldURL.String() {
		t.Errorf("Exepected (%q) to be rewritten. Got: %q", oldURL.String(), newURL.String())
	}
}
