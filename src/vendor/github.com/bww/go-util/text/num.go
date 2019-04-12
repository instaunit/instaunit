package text

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Errors
var ErrEmptyInput = fmt.Errorf("Empty input")

// Number parser function
type NumberParser func(string) (int, error)

// Range sorting
type byLowerBound [][]int

func (a byLowerBound) Len() int           { return len(a) }
func (a byLowerBound) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLowerBound) Less(i, j int) bool { return a[i][0] < a[j][0] }

// Parse a series of ranges from text. Multiple ranges can be
// represented in a comma-delimited list. Individual numbers in
// the list are treated as a range that includes only that number.
//
// An empty bound in a range is interpreted as "all the numbers
// from the constraining bound to the other bound". For example,
// the range "-10", with constraining bounds of [0, 100] means
// the range "0-10".
//
// Contiguous ranges are merged together.
//
// For example, the input: 0-5,5,6,7,8-10 produces: [[0,10]]
//
// And the input: 0-5,7,9-10 produces: [[0,5],[7,7],[9,10]]
//
func ParseRanges(s string, b []int) ([][]int, error) {
	return parseRanges(s, b, "-", strconv.Atoi)
}

// Parse ranges where the values are parsed by a function
func ParseRangesFunc(s string, b []int, d string, c NumberParser) ([][]int, error) {
	return parseRanges(s, b, d, c)
}

// Parse ranges
func parseRanges(s string, b []int, d string, c NumberParser) ([][]int, error) {
	if len(b) != 2 {
		return nil, fmt.Errorf("Invalid constraining bounds; length: %v", len(b))
	}
	if s == "" {
		return nil, ErrEmptyInput
	}

	var r [][]int
	for _, e := range strings.Split(s, ",") {
		var err error
		var l, u int

		e = strings.TrimSpace(e)
		if e == "" {
			return nil, fmt.Errorf("Invalid empty range")
		}

		if strings.Index(e, d) < 0 {
			if l, err = c(e); err != nil {
				return nil, fmt.Errorf("Invalid lower bound: %v", err)
			}
			u = l
		} else {
			n := strings.SplitN(e, d, 2)
			n[0], n[1] = strings.TrimSpace(n[0]), strings.TrimSpace(n[1])
			if n[0] == "" {
				l = b[0]
			} else if l, err = c(n[0]); err != nil {
				return nil, fmt.Errorf("Invalid lower bound: %v", err)
			}
			if n[1] == "" {
				u = b[1]
			} else if u, err = c(n[1]); err != nil {
				return nil, fmt.Errorf("Invalid upper bound: %v", err)
			}
		}

		r = append(r, []int{l, u})
	}

	sort.Sort(byLowerBound(r))

	i := 1
	for i < len(r) {
		e, p := r[i], r[i-1]
		if e[0] >= p[0] && e[0] <= p[1]+1 {
			if e[1] > p[1] {
				p[1] = e[1]
			} // extend right?
			r = append(r[:i], r[i+1:]...) // discard the extraneous range
		} else {
			i++
		}
	}

	return r, nil
}
