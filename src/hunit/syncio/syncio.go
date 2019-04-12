package syncio

import (
	"fmt"
	"io"
)

const buffer = 10

// Write lines in serial so multiple go routines don't interleave
// writes to the same writer
type syncWriter struct {
	io.Writer
	out chan []byte
}

func NewWriter(w io.Writer) *syncWriter {
	return (&syncWriter{w, make(chan []byte, buffer)}).Start()
}

func (w *syncWriter) Write(b []byte) (int, error) {
	w.out <- b
	return len(b), nil // assume success
}

func (w *syncWriter) Start() *syncWriter {
	go func() {
		for b := range w.out {
			_, err := w.Writer.Write(b)
			if err != nil {
				panic(fmt.Errorf("syncio: could not write: %v", err))
			}
		}
	}()
	return w
}
