package scan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const quote, esc = '\'', '\\'

// Test ident scanning
func TestScanIdent(t *testing.T) {
	assertScanIdent(t, `A, ok`, `A`, `, ok`, nil)
	assertScanIdent(t, `Hello123, ok`, `Hello123`, `, ok`, nil)
	assertScanIdent(t, `_Hello123, ok`, `_Hello123`, `, ok`, nil)
	assertScanIdent(t, `Hello_123, ok`, `Hello_123`, `, ok`, nil)
	assertScanIdent(t, `9Hello123, ok`, ``, ``, ErrInvalidSequence)
}

// Check an ident
func assertScanIdent(t *testing.T, in, ev, er string, eerr error) bool {
	av, ar, aerr := Ident(in)
	return assertScanFunc(t, in, ev, av, er, ar, eerr, aerr)
}

// Test string scanning
func TestScanString(t *testing.T) {
	assertScanString(t, `'Hello.' remainder`, `Hello.`, ` remainder`, nil)
	assertScanString(t, `'Hello. \a\f\n\r\t\v' remainder`, "Hello. \a\f\n\r\t\v", ` remainder`, nil)
	assertScanString(t, `'Hello, \'quoted\'' remainder`, `Hello, 'quoted'`, ` remainder`, nil)
	assertScanString(t, `'Hello, \\a\\f\\n\\r\\t\\v\\' remainder`, `Hello, \a\f\n\r\t\v\`, ` remainder`, nil)
	assertScanString(t, `'A \q' remainder`, ``, ``, ErrInvalidEscape) // invalid quote sequence
	assertScanString(t, `A' remainder`, ``, ``, ErrInvalidSequence)   // missing opening quote
	assertScanString(t, `'A remainder`, ``, ``, ErrInvalidSequence)   // missing closing quote
}

// Check a string
func assertScanString(t *testing.T, in, ev, er string, eerr error) bool {
	av, ar, aerr := String(in, quote, esc)
	return assertScanFunc(t, in, ev, av, er, ar, eerr, aerr)
}

// Check
func assertScanFunc(t *testing.T, in, ev, av, er, ar string, eerr, aerr error) bool {
	if aerr != nil || eerr != nil {
		fmt.Printf("%v -> %v\n", in, aerr)
		return assert.Equal(t, eerr, aerr, "Errors do not match")
	}
	fmt.Printf("%v -> [%v] [%v]\n", in, av, ar)
	res := true
	res = res && assert.Equal(t, ev, av, "Values do not match")
	res = res && assert.Equal(t, er, ar, "Remainders do not match")
	return res
}
