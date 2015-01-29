package heproxy

import (
	"encoding/xml"
	"io"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// Target represents a set of hostnames, potentially with wildcards.
type Target struct {
	Host string `xml:"host,attr"`
}

// Exclusion represents a URL pattern which should not be rewritten.
type Exclusion struct {
	Pattern   string `xml:"pattern,attr"`
	patternRE *pcre.Regexp
}

// Rule represents a rewritten URL.
type Rule struct {
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`

	fromRE *pcre.Regexp
}

// RuleSet represents a site, mapping from hostnames to HTTPS urls.
type RuleSet struct {
	Targets    []Target    `xml:"target"`
	Exclusions []Exclusion `xml:"exclusion"`
	Rules      []Rule      `xml:"rule"`
}

// ReadRuleSet returns a RuleSet read from `r`
func ReadRuleSet(r io.Reader) (*RuleSet, error) {
	decoder := xml.NewDecoder(r)
	ruleset := &RuleSet{}

	err := decoder.Decode(ruleset)
	return ruleset, err
}
