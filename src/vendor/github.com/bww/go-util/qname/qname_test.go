package util

import (
  "testing"
)

/**
 * Test qname
 */
func TestQName(t *testing.T) {
  
  // no options
  check(t, "*", "Anything!", 0, true)
  check(t, "*", "", 0, true)
  
  check(t, "a", "a", 0, true)
  check(t, "a", "b", 0, false)
  
  check(t, "a.b", "a.b", 0, true)
  check(t, "a.b", "a", 0, false)
  check(t, "a.b", "a.a", 0, false)
  
  check(t, "a.*", "a.a", 0, true)
  check(t, "a.*", "a.Z", 0, true)
  check(t, "a.*", "a.*", 0, true)
  
  check(t, "a.*", "a.b.c", 0, true)
  check(t, "a.*", "a.XYZ", 0, true)
  
  check(t, "a.*.z", "a.b", 0, false)
  check(t, "a.*.z", "a.z", 0, false)
  check(t, "a.*.z", "a.l.z", 0, true)
  
  // encompass
  check(t, "a", "a", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a", "b", MATCH_OPTION_ENCOMPASS, false)
  check(t, "a", "b.c", MATCH_OPTION_ENCOMPASS, false)
  
  check(t, "a.b", "a.b", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a.b", "a", MATCH_OPTION_ENCOMPASS, false)
  check(t, "a.b", "a.Z", MATCH_OPTION_ENCOMPASS, false)
  
  check(t, "a", "a.a", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a", "a.Z", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a", "a.*", MATCH_OPTION_ENCOMPASS, true)
  
  check(t, "a", "a.b.c", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a", "a.XYZ", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a.b", "a.b.c", MATCH_OPTION_ENCOMPASS, true)
  check(t, "a.b.c.d", "a.b.c", MATCH_OPTION_ENCOMPASS, false)
  
  check(t, "a.*.z", "a.b", MATCH_OPTION_ENCOMPASS, false)
  check(t, "a.*.z", "a.z", MATCH_OPTION_ENCOMPASS, false)
  check(t, "a.*.z", "a.l.z", MATCH_OPTION_ENCOMPASS, true)
  
}

/**
 * Check match
 */
func check(t *testing.T, l, r string, opts int, expect bool) {
  if QName(l).MatchesStringWithOptions(r, opts) != expect {
    if expect {
      t.Errorf("<%v> == <%v>", l, r)
    }else{
      t.Errorf("<%v> != <%v>", l, r)
    }
  }else{
    t.Logf("<%v> =~ <%v>", l, r)
  }
}
