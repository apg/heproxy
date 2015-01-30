package heproxy

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

// ErrPattern is returned as the error when a pattern is empty, but shouldn't be.
var ErrEmptyPattern = errors.New("Empty pattern")

// Target represents a set of hostnames, potentially with wildcards.
type Target struct {
	Host string `xml:"host,attr"`
}

// Exclusion represents a URL pattern which should not be rewritten.
type Exclusion struct {
	Pattern string

	patternRE pcre.Regexp
}

// Rule represents a rewritten URL.
type Rule struct {
	From string
	To   string

	fromRE pcre.Regexp
}

// RuleSet represents a site, mapping from hostnames to HTTPS urls.
type RuleSet struct {
	Name       string       `xml:"name,attr"`
	Targets    []*Target    `xml:"target"`
	Exclusions []*Exclusion `xml:"exclusion"`
	Rules      []*Rule      `xml:"rule"`
}

// NewRuleSet returns a RuleSet read from `r`
func NewRuleSet(r io.Reader) (*RuleSet, error) {
	decoder := xml.NewDecoder(r)
	ruleset := &RuleSet{}

	err := decoder.Decode(ruleset)
	return ruleset, err
}

// Includes tests to see if url is included by this RuleSet
func (r *RuleSet) Includes(u *url.URL) bool {
	if !r.Excludes(u) {
		for _, t := range r.Targets {
			if t.matches(u) {
				return true
			}
		}
	}
	return false
}

// Excludes tests to see if url is excluded by this Ruleset
func (r *RuleSet) Excludes(url *url.URL) bool {
	for _, e := range r.Exclusions {
		if e.matches(url) {
			return true
		}
	}
	return false
}

// Rewrite attempts to rewrite the url.
func (r *RuleSet) Rewrite(url *url.URL) *url.URL {
	for _, rule := range r.Rules {
		if nURL, success := rule.Apply(url); success {
			return nURL
		}
	}

	return url
}

// UnmarshalXML unmarshals and compiles a regular expression pulled out of the XML.
func (r *Rule) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	var dummy interface{}           // consume the rest of the element.
	d.DecodeElement(&dummy, &start) // not sure about this.
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "from":
			r.From = attr.Value
			r.fromRE, err = compile(r.From)
			if err != nil {
				return err
			}
		case "to":
			r.To = attr.Value
		}
	}

	return err
}

// UnmarshalXML unmarshals and compiles a regular expression pulled out of the XML.
func (e *Exclusion) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	var dummy interface{}           // consume the rest of the element.
	d.DecodeElement(&dummy, &start) // not sure about this.
	for _, attr := range start.Attr {
		if attr.Name.Local == "pattern" {
			e.Pattern = attr.Value
			e.patternRE, err = compile(e.Pattern)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Apply attempts to apply rule to given url
func (r *Rule) Apply(old *url.URL) (new *url.URL, found bool) {
	var err error

	oldStr := old.String()
	matcher := r.fromRE.MatcherString(oldStr, 0)
	if matcher.Matches() {
		newStr := r.To
		for i := 0; i < matcher.Groups(); i++ {
			newStr = strings.Replace(newStr, fmt.Sprintf("$%d", i+1), string(matcher.Group(i+1)), -1)
		}

		// replace oldStr's matched part, with
		fullStr := strings.Replace(oldStr, string(matcher.Group(0)), newStr, 1)
		new, err = url.Parse(fullStr)
		if err != nil {
			log.Printf("Fudge: %q, when parsing: %q", err, fullStr)
		}
		return new, true
	}

	return old, false
}

func (e *Exclusion) matches(u *url.URL) bool {
	if e.patternRE.MatcherString(u.String(), 0).Matches() {
		return true
	}
	return false
}

func (t *Target) matches(u *url.URL) bool {
	// FIXME: this doesn't do ports or anything like that.
	return u.Host == t.Host
}

func compile(pat string) (pcre.Regexp, error) {
	re, cErr := pcre.Compile(massagePattern(pat), 0)
	if cErr != nil {
		return re, fmt.Errorf(cErr.String())
	}
	return re, nil
}

// may not be necessary, but for now...
func massagePattern(pat string) string {
	return pat
}
