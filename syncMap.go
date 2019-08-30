package main

import (
	"os/exec"
	"sync"
	"time"
)

// SyncMap is a map of taskID to previous processes.
type SyncMap struct {
	cmdMap map[string]*CommandWrapper
	*sync.Mutex
}

// NewMap is used to return an empty SyncMap.
func NewMap() SyncMap {
	return SyncMap{
		cmdMap: make(map[string]*CommandWrapper),
		Mutex:  &sync.Mutex{},
	}
}

// Put is an operation insert a new process in a SyncMap.
func (rwm *SyncMap) Put(taskID string, cmd *CommandWrapper) {
	rwm.Lock()
	defer rwm.Unlock()
	rwm.cmdMap[taskID] = cmd
}

// Get returns the CommandWrapper mapped to the taskID.
func (rwm *SyncMap) Get(taskID string) (cmd *CommandWrapper, ok bool) {
	rwm.Lock()
	defer rwm.Unlock()
	cmd, ok = rwm.cmdMap[taskID]
	return cmd, ok
}

// CommandWrapper encapsulates an exec.Cmd object with status metadata.
type CommandWrapper struct {
	Command    *exec.Cmd // underlying command
	StartTime  time.Time
	endTime    *time.Time
	StdoutBuff *Output
	execError  *string
	exitCode   int
	*sync.Mutex
}

// NewCommandWrapper is used to return a default CommandWrapper object.
func NewCommandWrapper(cmd *exec.Cmd, outBuff *Output) *CommandWrapper {
	return &CommandWrapper{
		Command:    cmd,
		StartTime:  time.Now(),
		endTime:    nil,
		StdoutBuff: outBuff,
		exitCode:   cmd.ProcessState.ExitCode(),
		execError:  nil,
		Mutex:      &sync.Mutex{},
	}
}

// SetEndTime is used to set the end time of a process.
func (cw *CommandWrapper) SetEndTime(endTime time.Time) {
	cw.Lock()
	defer cw.Unlock()
	cw.endTime = &endTime
}

// GetEndTime is used to access the end time of a process.
func (cw *CommandWrapper) GetEndTime() *time.Time {
	cw.Lock()
	defer cw.Unlock()
	return cw.endTime
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

// GetExecError returns errors attributed to the CommandWrapper lifecycle.
func (cw *CommandWrapper) GetExecError() *string {
	cw.Lock()
	defer cw.Unlock()
	return cw.execError
}

// SetExecError appends an error attributed to the CommandWrapper lifecycle.
func (cw *CommandWrapper) SetExecError(error string) {
	cw.Lock()
	defer cw.Unlock()
	cw.execError = &error
}
