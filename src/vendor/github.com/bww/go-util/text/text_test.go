package text

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

/**
 * Test normalize string
 */
func TestNormalizeString(t *testing.T) {
  assert.Equal(t, "", NormalizeString(" "))
  assert.Equal(t, "", NormalizeString("  "))
  assert.Equal(t, "a", NormalizeString(" a "))
  assert.Equal(t, "abc", NormalizeString(" ABC "))
  assert.Equal(t, "abc a", NormalizeString(" ABC A"))
  assert.Equal(t, "abc a", NormalizeString(" ABC   A "))
  assert.Equal(t, "a,b,c", NormalizeString(" A,B,C "))
  assert.Equal(t, "a, b, c", NormalizeString(" A,  B,  C "))
  assert.Equal(t, "a_b_c_a", NormalizeString(" A_B_C_A "))
  assert.Equal(t, "a-b-c-a", NormalizeString(" A-B-C-A "))
  assert.Equal(t, "abc 123", NormalizeString("ABC 123"))
  assert.Equal(t, "abc-123", NormalizeString("ABC-123"))
  assert.Equal(t, "abc 123", NormalizeString("ABC/123"))
  assert.Equal(t, "abc 123", NormalizeString("ABC///123"))
  assert.Equal(t, "ça c'est bien passée", NormalizeString("Ça c'est bien passée   !"))
  assert.Equal(t, "abc", NormalizeString(`
  abc`))
}

/**
 * Test normalize terms
 */
func TestNormalizeTerms(t *testing.T) {
  assert.Equal(t, "", NormalizeTerms(" "))
  assert.Equal(t, "", NormalizeTerms("  "))
  assert.Equal(t, "a", NormalizeTerms(" a "))
  assert.Equal(t, "abc", NormalizeTerms(" ABC "))
  assert.Equal(t, "abc a", NormalizeTerms(" ABC A"))
  assert.Equal(t, "abc a", NormalizeTerms(" ABC   A "))
  assert.Equal(t, "a_b_c_a", NormalizeTerms(" A_B_C_A "))
  assert.Equal(t, "a-b-c-a", NormalizeTerms(" A-B-C-A "))
  assert.Equal(t, "abc 123", NormalizeTerms("ABC 123"))
  assert.Equal(t, "abc-123", NormalizeTerms("ABC-123"))
  assert.Equal(t, "abc 123", NormalizeTerms("ABC/123"))
  assert.Equal(t, "abc 123", NormalizeTerms("ABC///123"))
  assert.Equal(t, "ça c'est bien passée", NormalizeTerms("Ça c'est bien passée   !"))
  assert.Equal(t, "abc", NormalizeTerms(`
  abc`))
}

/**
 * Test collapse spaces
 */
func TestCollapseSpaces(t *testing.T) {
  assert.Equal(t, "", CollapseSpaces(""))
  assert.Equal(t, "", CollapseSpaces(" "))
  assert.Equal(t, "", CollapseSpaces("   "))
  assert.Equal(t, "a", CollapseSpaces(" a "))
  assert.Equal(t, "A B C", CollapseSpaces(" A B C "))
  assert.Equal(t, "A$B$C", CollapseSpaces(" A$B$C "))
}

/**
 * Test normalize join
 */
func TestNormalizeJoin(t *testing.T) {
  assert.Equal(t, "", NormalizeJoin([]string{}, "a", "b"))
  assert.Equal(t, "1", NormalizeJoin([]string{"1"}, ", ", " and "))
  assert.Equal(t, "1 and 2", NormalizeJoin([]string{"1", "2"}, ", ", " and "))
  assert.Equal(t, "1, 2 and 3", NormalizeJoin([]string{"1", "2", "3"}, ", ", " and "))
  assert.Equal(t, "1, 2, 3 and 4", NormalizeJoin([]string{"1", "2", "3", "4"}, ", ", " and "))
}
