package heproxy

import "net/url"

// Rewriter holds rulesets for remapping URLs
type Rewriter struct {
	rulesets []*RuleSet
}

// NewRewriter returns a Remapper which can remap URLs based on the RuleSets
func NewRewriter(rs []*RuleSet) *Rewriter {
	return &Rewriter{
		rulesets: rs,
	}
}

// Rewrite returns a new URL based on the rulesets
func (r *Rewriter) Rewrite(u *url.URL) *url.URL {
	// FIXME: Brute force approach here, which obviously won't scale.

	// Find a possible ruleset that matches this URL.
	for _, r := range r.rulesets {
		if r.Includes(u) {
			return r.Rewrite(u)
		}
	}

	return u
}
