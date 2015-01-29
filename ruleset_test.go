package heproxy

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRuleSetReader(t *testing.T) {
	set := `<ruleset name="EFF">
  <target host="eff.org" />
  <target host="*.eff.org" />

  <exclusion pattern="^http://action\.eff\.org/"/>
  <rule from="^http://eff\.org/" to="https://eff.org/"/>
  <rule from="^http://([^/:@]*)\.eff\.org/" to="https://$1.eff.org/"/>
</ruleset>`

	r := bytes.NewBufferString(set)

	ruleSet, err := ReadRuleSet(r)
	fmt.Printf("%v :::: %s", ruleSet, err)
}
