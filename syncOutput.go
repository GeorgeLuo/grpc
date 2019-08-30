package main

import (
	"bufio"
	"bytes"
	"sync"
)

// Output is used to retrieve output to a buffer as populated.
type Output struct {
	buf *bytes.Buffer
	*sync.Mutex
}

// NewOutput is used to return a default Output object.
func NewOutput() *Output {
	return &Output{
		buf:   &bytes.Buffer{},
		Mutex: &sync.Mutex{},
	}
}

// Write is an operation to write to the underlying buffer.
func (rw *Output) Write(p []byte) (int, error) {
	rw.Lock()
	defer rw.Unlock()
	return rw.buf.Write(p)
}

// Lines returns the contents of the buffer at the current state.
func (rw *Output) Lines() []string {
	rw.Lock()
	defer rw.Unlock()
	s := bufio.NewScanner(rw.buf)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines
}
