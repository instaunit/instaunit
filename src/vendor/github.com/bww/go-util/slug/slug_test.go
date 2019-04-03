package slug

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/**
 * Test slugs
 */
func TestSlugs(t *testing.T) {
	assert.Equal(t, `this-is-a-title-to-slugify`, Slugify(`This is a Title To Slugify`))
	assert.Equal(t, `this-is-a-title-to-slugify`, Slugify(`  This is a   Title   To Slugify  `))
	assert.Equal(t, `this-is-a-title-to-slugify`, Slugify(`This is a &quot;Title&quot; To Slugify`))
	assert.Equal(t, `this-and-that`, Slugify(`This &amp; That`))
	assert.Equal(t, `the-mayor-said-no-way`, Slugify(`The Mayor Said «No Way!»`))
	assert.Equal(t, `pay-back-the-100-dollars-25-cents-you-owe`, Slugify(`Pay back the $100.25 you owe!`))
	assert.Equal(t, ``, Slugify(`          `))
	assert.Equal(t, ``, Slugify(` `))
	assert.Equal(t, ``, Slugify(``))
	assert.Equal(t, `this-already-has-hyphens`, Slugify(`THIS-already-HAS-hyphens`))
	assert.Equal(t, `this-already-has-too-many-hyphens`, Slugify(`this---already---has---too---many---hyphens`))
	assert.Equal(t, ``, Slugify(`----------`))
	assert.Equal(t, ``, Slugify(`-`))
	assert.Equal(t, ``, Slugify(` - `))
	assert.Equal(t, ``, Slugify(` - - - `))
	assert.Equal(t, `a`, Slugify(` a `))
	assert.Equal(t, `a-b`, Slugify(` a-b `))
	assert.Equal(t, `a-b`, Slugify(` a---b `))
}

/**
 * Test diacritics
 */
func TestNormalizeDiacritics(t *testing.T) {
	assert.Equal(t, "aaeeiiooouuuyn", normalizeDiacritics("áàéèíìóòöúùüÿñ"))
	assert.Equal(t, "ca", normalizeDiacritics("ça"))
	assert.Equal(t, "sur", normalizeDiacritics("sûr"))
	assert.Equal(t, "Wolter", normalizeDiacritics("Wölter"))
}

/**
 * Test smart punctuation
 */
func TestNormalizeSmartPunctuation(t *testing.T) {
	assert.Equal(t, `""""""""""""`, normalizeSmartPunctuation(`«»“”„‟❝❞〝〞〟＂`))
	assert.Equal(t, `''''''''`, normalizeSmartPunctuation(`‘’‚‛‹›❛❜`))
	assert.Equal(t, `...`, normalizeSmartPunctuation(`…`))
}

/**
 * Test strip control characters
 */
func TestStripControlCharacters(t *testing.T) {
	s := "a"
	for i := 0; i < 20; i++ {
		s += string(i)
	}
	s += "b"
	assert.Equal(t, "ab", stripControlCharacters(s))
}

/**
 * Test convert symbols to words
 */
func TestConvertSymbolsToWords(t *testing.T) {
	assert.Equal(t, `This and That`, convertSymbolsToWords(`This & That`))
	assert.Equal(t, `He's at the movies`, convertSymbolsToWords(`He's @ the movies`))
	assert.Equal(t, `98 degrees `, convertSymbolsToWords(`98º`))
	assert.Equal(t, `98 divided by 100`, convertSymbolsToWords(`98 ÷ 100`))
	assert.Equal(t, `101 dot 1 FM`, convertSymbolsToWords(`101.1 FM`))
	assert.Equal(t, `Nice job ellipsis `, convertSymbolsToWords(`Nice job...`))
	assert.Equal(t, `100 equals 2 times 50`, convertSymbolsToWords(`100 = 2 times 50`))
	assert.Equal(t, `99 percent `, convertSymbolsToWords(`99%`))
	assert.Equal(t, `He slash she`, convertSymbolsToWords(`He/she`))
	assert.Equal(t, `Baby, you're a star `, convertSymbolsToWords(`Baby, you're a *`))
}

/**
 * Test currnecy symbols to words
 */
func TestConvertCurrenciesToWords(t *testing.T) {
	assert.Equal(t, `He paid 100 dollars back`, convertCurrenciesToWords(`He paid $100 back`))
	assert.Equal(t, `He paid 100 dollars 10 cents back`, convertCurrenciesToWords(`He paid $100.10 back`))
	assert.Equal(t, `He paid 100 euros back`, convertCurrenciesToWords(`He paid €100 back`))
	assert.Equal(t, `He paid 100 euros 10 cents back`, convertCurrenciesToWords(`He paid €100.10 back`))
	assert.Equal(t, `He paid 100 pounds back`, convertCurrenciesToWords(`He paid £100 back`))
	assert.Equal(t, `He paid 100 pounds 10 pence back`, convertCurrenciesToWords(`He paid £100.10 back`))
	assert.Equal(t, `He paid 100 reais back`, convertCurrenciesToWords(`He paid R$100 back`))
	assert.Equal(t, `He paid 100 reais 10 cents back`, convertCurrenciesToWords(`He paid R$100.10 back`))
}
