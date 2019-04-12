package text

import (
	"fmt"
	"io"
)

// Write a string in the manner of the output of `hexdump -C`,
// with byte values on the left corresponding to characters
// on the right, when printable
func Hexdump(w io.Writer, b []byte, cols int) {
	f := fmt.Sprintf("%%-%ds | %%s", (cols*3)-1)

	var i int
	var e byte
	if cols < 1 {
		cols = 1
	}

	var cp, ch string
	for i, e = range b {
		if i > 0 && (i%cols) == 0 {
			if i > cols {
				fmt.Fprintln(w)
			}
			fmt.Fprintf(w, f, cp, ch)
			cp, ch = "", ""
		} else if i != 0 {
			cp += " "
		}
		cp += fmt.Sprintf("%02x", int(e))
		if e >= 0x20 && e < 0x7f {
			ch += string(e)
		} else {
			ch += "."
		}
	}

	if len(cp) > 0 {
		if i > cols {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, f, cp, ch)
	}
}
