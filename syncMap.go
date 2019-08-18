package main

import (
	"os/exec"
	"sync"
	"time"
)

// Output is used to retrieve output to a buffer as populated.
type SyncMap struct {
	cmdMap map[string]*CommandWrapper
	*sync.Mutex
}

// NewOutput is used to return a default Output object.
func NewMap() SyncMap {
	return SyncMap{
		cmdMap: make(map[string]*CommandWrapper),
		Mutex:  &sync.Mutex{},
	}
}

// Write is an operation to write to the underlying buffer.
func (rwm *SyncMap) Put(task_id string, cmd *CommandWrapper) {
	rwm.Lock()
	defer rwm.Unlock()
	rwm.cmdMap[task_id] = cmd
}

// Lines returns the contents of the buffer at the current state.
func (rwm *SyncMap) Get(task_id string) (cmd *CommandWrapper, ok bool) {
	rwm.Lock()
	defer rwm.Unlock()
	cmd, ok = rwm.cmdMap[task_id]
	return cmd, ok
}

// CommandWrapper encapsulates an exec.Cmd object with status metadata.
type CommandWrapper struct {
	Command    *exec.Cmd // underlying command
	finished   bool      // set upon process finish
	StartTime  time.Time
	EndTime    time.Time
	StdoutBuff *Output
	exitCode   int
	*sync.Mutex
}

// NewCommandWrapper is used to return a default CommandWrapper object.
func NewCommandWrapper(cmd *exec.Cmd, outBuff *Output) *CommandWrapper {
	return &CommandWrapper{
		Command:    cmd,
		finished:   false,
		StartTime:  time.Now(),
		EndTime:    time.Time{},
		StdoutBuff: outBuff,
		exitCode:   cmd.ProcessState.ExitCode(),
		Mutex:      &sync.Mutex{},
	}
}

// SetFinished is used to set the finished status of a process.
func (cw *CommandWrapper) SetFinished(finished bool) {
	cw.Lock()
	defer cw.Unlock()
	cw.finished = finished
}

// GetFinished is used to access the finished status of a process.
func (cw *CommandWrapper) GetFinished() bool {
	cw.Lock()
	defer cw.Unlock()
	return cw.finished
}

// GetExitCode is used to access exitCode value.
func (cw *CommandWrapper) GetExitCode() int {
	cw.Lock()
	defer cw.Unlock()
	return cw.exitCode
}

// SetExitCode is used set the exitCode value.
func (cw *CommandWrapper) SetExitCode(code int) {
	cw.Lock()
	defer cw.Unlock()
	cw.exitCode = code
}
