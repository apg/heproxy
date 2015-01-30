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
	Pattern   string `xml:"pattern,attr"`
	patternRE pcre.Regexp
}

// Rule represents a rewritten URL.
type Rule struct {
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`

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
	if err == nil {
		err = ruleset.compile()
	}
	return ruleset, err
}

// Matches tests to see if url is covered by this RuleSet
func (r *RuleSet) Matches(url *url.URL) bool {
	for _, t := range r.Targets {
		// FIXME: refactor into a Target.Matches(), to take into consideration more things.
		if t.Host == url.Host { // FIXME: url.Host might have a port, too
			// unless it's excluded...
			for _, e := range r.Exclusions {
				if e.patternRE.MatcherString(url.String(), 0).Matches() {
					return false
				}
			}
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

func (r *RuleSet) compile() error {
	var err error

	for _, t := range r.Targets {
		err = t.compile()
		if err != nil {
			return err
		}
	}

	for _, t := range r.Exclusions {
		err = t.compile()
		if err != nil {
			return err
		}
	}

	for _, t := range r.Rules {
		err = t.compile()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Target) compile() error {
	return nil
}

func (e *Exclusion) compile() error {
	var err error

	if e.Pattern != "" {
		e.patternRE, err = compile(e.Pattern)
	} else {
		err = ErrEmptyPattern
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

func (r *Rule) compile() error {
	var err error
	if r.From != "" {
		r.fromRE, err = compile(r.From)
	} else {
		err = ErrEmptyPattern
	}
	return err
}

func compile(pat string) (pcre.Regexp, error) {
	re, cErr := pcre.Compile(massagePattern(pat), 0)
	if cErr != nil {
		return re, errors.New(cErr.String())
	}
	return re, nil
}

// may not be necessary
func massagePattern(pat string) string {
	return pat
}
