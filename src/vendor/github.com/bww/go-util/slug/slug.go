package slug

import (
  "html"
  "regexp"
  "unicode"
  "golang.org/x/text/transform"
  "golang.org/x/text/unicode/norm"
)

// Slugify a string. We attempt to do this so as to produce slugs
// that are as identical as possible to the (Ruby) Mongoid slug
// package. That package actually uses another package, Stringex
// (https://github.com/rsl/stringex) to generate the slug strings.
// Stringex, in turn, provides a special adapater for Mongoid. We
// use Stringex as the reference implementation without the
// exclusion, limit or downcase options.
// 
//    def to_url(options = {})
//      return self if options[:exclude] && options[:exclude].include?(self)
//      options = stringex_default_options.merge(options)
//      whitespace_replacement_token = options[:replace_whitespace_with]
//      dummy = remove_formatting(options).
//                replace_whitespace(whitespace_replacement_token).
//                collapse(whitespace_replacement_token).
//                limit(options[:limit], options[:truncate_words], whitespace_replacement_token)
//      dummy.downcase! unless options[:force_downcase] == false
//      dummy
//    end
// 
// This method makes use of the following routine to normalize all
// kinds of things. We don't perform all of these normalizations.
// 
//    def remove_formatting(options = {})
//      strip_html_tags.
//        convert_smart_punctuation.
//        convert_accented_html_entities.
//        convert_vulgar_fractions.
//        convert_unreadable_control_characters.
//        convert_miscellaneous_html_entities.
//        convert_miscellaneous_characters(options).
//        to_ascii.
//        # NOTE: String#to_ascii may convert some Unicode characters to ascii we'd already transliterated
//        # so we need to do it again just to be safe
//        convert_miscellaneous_characters(options).
//        collapse
//    end
// 
func Slugify(s string) string {
  return slugify(s, "-")
}

// Slugify
func slugify(s, w string) string {
  var g string
  
  s = StripHTMLTags(s)
  s = html.UnescapeString(s)
  s = stripControlCharacters(s)
  s = normalizeDiacritics(s)
  s = convertCurrenciesToWords(s)
  s = convertSymbolsToWords(s)
  
  sp := 0
  for _, e := range s {
    if unicode.IsSpace(e) {
      sp++
      continue
    }
    if unicode.IsLetter(e) || unicode.IsNumber(e) {
      if sp > 0 {
        // add the space if we're not at the beginning or the end
        if len(g) > 0 {
          g += w
        }
        sp = 0 // clear space
      }
      g += string(unicode.ToLower(e))
    }
  }
  
  return g
}

// Non-space marks
func isMn(r rune) bool {
  return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// Convert characters with diacritical marks to their unaccented/base
// counterparts. See also:
// http://stackoverflow.com/questions/26722450/remove-diacritics-using-go
func normalizeDiacritics(s string) string {
  t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
  result, _, _ := transform.String(t, s)
  return result
}

// Map smart punctuation to it's dumb counterpart
var smartPunct = map[rune]string {
  '«': "\"", '»': "\"",
  '“': "\"", '”': "\"",
  '„': "\"", '‟': "\"",
  '❝': "\"", '❞': "\"",
  '〝': "\"", '〞': "\"",
  '〟': "\"", '＂': "\"",
  '‘': "'",  '’': "'",
  '‚': "'",  '‛': "'",
  '‹': "'",  '›': "'",
  '❛': "'",  '❜': "'",
  '…': "...",
}

// Convert special quotation marks to ASCII " runes.
func normalizeSmartPunctuation(s string) string {
  var c string
  for _, e := range s {
    if r, ok := smartPunct[e]; ok {
      c += r
    }else{
      c += string(e)
    }
  }
  return c
}

// Strip out control characters
func stripControlCharacters(s string) string {
  var c string
  for _, e := range s {
    if !unicode.IsControl(e) {
      c += string(e)
    }
  }
  return c
}

// Expression / replacement
type regexpReplace struct {
  Expr    *regexp.Regexp
  Replace string
}

// Symbols for term conversion (this is a pretty expensive way to do this...)
var termSymbols = []regexpReplace {
  regexpReplace{regexp.MustCompile(`\s*&\s*`), " and "},
  regexpReplace{regexp.MustCompile(`\s*@\s*`), " at "},
  regexpReplace{regexp.MustCompile(`\s*º\s*`), " degrees "},
  regexpReplace{regexp.MustCompile(`\s*°\s*`), " degrees "},
  regexpReplace{regexp.MustCompile(`\s*÷\s*`), " divided by "},
  regexpReplace{regexp.MustCompile(`\s*\.{3,}\s*`), " ellipsis "},
  regexpReplace{regexp.MustCompile(`(\S|^)\.(\S)`), "$1 dot $2"},
  regexpReplace{regexp.MustCompile(`\s*=\s*`), " equals "},
  regexpReplace{regexp.MustCompile(`\s*%\s*`), " percent "},
  regexpReplace{regexp.MustCompile(`\s*(\\|\/|／)\s*`), " slash "},
  regexpReplace{regexp.MustCompile(`\s*\*\s*`), " star "},
}

// Convert certain special symbols to their word counterpard
func convertSymbolsToWords(s string) string {
  for _, e := range termSymbols {
    s = e.Expr.ReplaceAllString(s, e.Replace)
  }
  return s
}

// Symbols for currency conversion (this is a pretty expensive way to do this...)
var currencySymbols = []regexpReplace {
  regexpReplace{regexp.MustCompile(`(?:\s|^)€(\d+)(?:\s|$)`), " $1 euros "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)€(\d+)\.(\d+)(?:\s|$)`), " $1 euros $2 cents "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)\$(\d+)(?:\s|$)`), " $1 dollars "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)\$(\d+)\.(\d+)(?:\s|$)`), " $1 dollars $2 cents "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)£(\d+)(?:\s|$)`), " $1 pounds "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)£(\d+)\.(\d+)(?:\s|$)`), " $1 pounds $2 pence "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)¥(\d+)(?:\s|$)`), " $1 yen "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)R\$(\d+)(?:\s|$)`), " $1 reais "},
  regexpReplace{regexp.MustCompile(`(?:\s|^)R\$(\d+)\.(\d+)(?:\s|$)`), " $1 reais $2 cents "},
}

// Convert certain special symbols to their word counterpard
func convertCurrenciesToWords(s string) string {
  for _, e := range currencySymbols {
    s = e.Expr.ReplaceAllString(s, e.Replace)
  }
  return s
}

