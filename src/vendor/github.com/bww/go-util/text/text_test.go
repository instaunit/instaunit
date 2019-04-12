package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestIdentizeString(t *testing.T) {
	assert.Equal(t, "", IdentizeString(" "))
	assert.Equal(t, "", IdentizeString("  "))
	assert.Equal(t, "a", IdentizeString(" a "))
	assert.Equal(t, "abc", IdentizeString(" ABC "))
	assert.Equal(t, "abc_a", IdentizeString(" ABC A"))
	assert.Equal(t, "abc_a", IdentizeString(" ABC   A "))
	assert.Equal(t, "a_b_c_a", IdentizeString(" A_B_C_A "))
	assert.Equal(t, "a_b_c_a", IdentizeString(" A-B-C-A "))
	assert.Equal(t, "abc_123", IdentizeString("ABC 123"))
	assert.Equal(t, "abc_123", IdentizeString("ABC-123"))
	assert.Equal(t, "abc_123", IdentizeString("ABC/123"))
	assert.Equal(t, "abc_123", IdentizeString("ABC///123"))
	assert.Equal(t, "ça_c_est_bien_passée", IdentizeString("Ça c'est bien passée   !"))
	assert.Equal(t, "abc", IdentizeString(`
  abc`))
}

func TestCollapseSpaces(t *testing.T) {
	assert.Equal(t, "", CollapseSpaces(""))
	assert.Equal(t, "", CollapseSpaces(" "))
	assert.Equal(t, "", CollapseSpaces("   "))
	assert.Equal(t, "a", CollapseSpaces(" a "))
	assert.Equal(t, "A B C", CollapseSpaces(" A B C "))
	assert.Equal(t, "A$B$C", CollapseSpaces(" A$B$C "))
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "", Truncate("", 10, ""))
	assert.Equal(t, "Hello", Truncate("Hello", 10, "..."))
	assert.Equal(t, "Hello dude", Truncate("Hello dude, how are you?", 10, ""))
	assert.Equal(t, "Hello dude...", Truncate("Hello dude, how are you?", 10, "..."))
	assert.Equal(t, "Hello 日本, ", Truncate("Hello 日本, how are you?", 10, ""))
	assert.Equal(t, "Hello 日本, ...", Truncate("Hello 日本, how are you?", 10, "..."))
}

func TestNormalizeJoin(t *testing.T) {
	assert.Equal(t, "", NormalizeJoin([]string{}, "a", "b"))
	assert.Equal(t, "1", NormalizeJoin([]string{"1"}, ", ", " and "))
	assert.Equal(t, "1 and 2", NormalizeJoin([]string{"1", "2"}, ", ", " and "))
	assert.Equal(t, "1, 2 and 3", NormalizeJoin([]string{"1", "2", "3"}, ", ", " and "))
	assert.Equal(t, "1, 2, 3 and 4", NormalizeJoin([]string{"1", "2", "3", "4"}, ", ", " and "))
}

func TestCoalesce(t *testing.T) {
	assert.Equal(t, "", Coalesce())
	assert.Equal(t, "", Coalesce(""))
	assert.Equal(t, "", Coalesce("", ""))
	assert.Equal(t, "a", Coalesce("", "a"))
	assert.Equal(t, "a", Coalesce("a", "b"))
}

func TestStripControl(t *testing.T) {
	assert.Equal(t, "Hello", StripControl("Hello"))
	assert.Equal(t, "Hello", StripControl("Hel\u0000lo"))
	assert.Equal(t, "Hello", StripControl("Hello\u0000"))
	assert.Equal(t, "Hello", StripControl("\u0000Hello"))
	assert.Equal(t, "Hello", StripControl("\u0000Hello\u0000"))
}

func TestNormalizeDiacritics(t *testing.T) {
	assert.Equal(t, "Ca c'est bien passee!", NormalizeDiacritics("Ça c'est bien passée!"))
	assert.Equal(t, "We love umlauts!", NormalizeDiacritics("We love ümlauts!"))
	assert.Equal(t, "aaaeeeiiiooouuuu", NormalizeDiacritics("âàáêèéîìíôòóüûùú"))
	assert.Equal(t, "AAAEEEIIIOOOUUUU", NormalizeDiacritics("ÂÀÁÊÈÉÎÌÍÔÒÓÜÛÙÚ"))
}
