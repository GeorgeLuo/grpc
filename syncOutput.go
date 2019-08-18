package main

import (
	"bufio"
	"bytes"
	"sync"
)

// Output is used to retrieve output to a buffer as populated.
type Output struct {
	buf   *bytes.Buffer
	lines []string
	*sync.Mutex
}

// NewOutput is used to return a default Output object.
func NewOutput() *Output {
	return &Output{
		buf:   &bytes.Buffer{},
		lines: []string{},
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
	for s.Scan() {
		rw.lines = append(rw.lines, s.Text())
	}
	return rw.lines
}
