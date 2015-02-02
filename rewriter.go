package heproxy

import (
	"net/url"
	"strings"
)

// Rewriter rewrites url to a new url
type Rewriter interface {
	Rewrite(*url.URL) *url.URL
}

// LinearRewriter holds rulesets for remapping URLs
type LinearRewriter struct {
	rulesets []*RuleSet
}

// NewLinearRewriter returns a Remapper which can remap URLs based on the RuleSets
func NewLinearRewriter(rs []*RuleSet) Rewriter {
	return &LinearRewriter{
		rulesets: rs,
	}
}

// Rewrite returns a new URL based on the contained rulesets
func (r *LinearRewriter) Rewrite(u *url.URL) *url.URL {
	// Find a possible ruleset that matches this URL.
	for _, r := range r.rulesets {
		if r.Includes(u) {
			return r.Rewrite(u)
		}
	}

	return u
}

type trieNode struct {
	// If this is nil, it's an internal node.
	Rset *RuleSet

	children map[string]*trieNode
}

// TrieRewriter utilizes a trie tree for rewriting urls based on rulesets
type TrieRewriter struct {
	root *trieNode
}

// NewTrieRewriter returns a TrieRewriter for remapping urls.
func NewTrieRewriter(rss []*RuleSet) Rewriter {
	t := new(TrieRewriter)
	t.build(rss)
	return t
}

// Rewrite returns a new URL based on the contained rulesets
func (t *TrieRewriter) Rewrite(u *url.URL) *url.URL {
	host := strings.Split(u.Host, ":")[0]
	bits := strings.Split(host, ".")

	rs := t.root.findMatch(bits)
	if rs != nil && rs.Includes(u) {
		return rs.Rewrite(u)
	}
	return u
}

func (t *TrieRewriter) build(rss []*RuleSet) {
	if t.root == nil {
		t.root = &trieNode{children: make(map[string]*trieNode)}
	}

	for _, rs := range rss {
		for _, target := range rs.Targets {
			bits := strings.Split(target.Host, ".")
			t.root.add(bits, rs)
		}
	}
}

func (t *trieNode) add(pieces []string, rs *RuleSet) {
	if len(pieces) == 0 {
		t.Rset = rs
		return
	}

	if next, ok := t.children[pieces[0]]; !ok {
		// add to current children, and recurse
		next = &trieNode{children: make(map[string]*trieNode)}
		t.children[pieces[0]] = next
		next.add(pieces[1:], rs)
	} else {
		next.add(pieces[1:], rs)
	}
}

func (t *trieNode) findMatch(target []string) *RuleSet {
	// step through target. If t.children has a '*', but no direct match, explore each child searching for a match.

	switch len(target) {
	case 0:
		return t.Rset

	case 1:
		if node, ok := t.children[target[0]]; !ok {
			// is there a wild? If so, return the Rset
			if wild, ok := t.children["*"]; ok {
				return wild.findMatch(target[1:])
			}
		} else {
			return node.findMatch(target[1:])
		}
	default:
		if node, ok := t.children[target[0]]; !ok {
			if wild, ok := t.children["*"]; ok {
				return wild.findMatch(target[1:])
			}
		} else {
			return node.findMatch(target[1:])
		}
	}

	return nil
}
