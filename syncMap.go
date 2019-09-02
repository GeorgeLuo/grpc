package main

import (
	"os/exec"
	"sync"
	"time"
)

// SyncMap is a map of taskID to previous processes. The abstraction is defined
// to protect enforce thread-safety of read and write to CommandWrapper objects.
type SyncMap struct {
	cmdMap map[string]*CommandWrapper
	mutex  *sync.Mutex
}

// NewMap is used to return an empty SyncMap.
func NewMap() SyncMap {
	return SyncMap{
		cmdMap: make(map[string]*CommandWrapper),
		mutex:  &sync.Mutex{},
	}
}

// Put is an operation insert a new process in a SyncMap.
func (rwm *SyncMap) Put(taskID string, cmd *CommandWrapper) {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()
	rwm.cmdMap[taskID] = cmd
}

// Get returns the CommandWrapper mapped to the taskID.
func (rwm *SyncMap) Get(taskID string) (cmd *CommandWrapper, ok bool) {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()
	cmd, ok = rwm.cmdMap[taskID]
	return cmd, ok
}

// CommandWrapper encapsulates an exec.Cmd object with status metadata.
type CommandWrapper struct {
	Command    *exec.Cmd // underlying command
	StartTime  time.Time
	endTime    *time.Time
	StdoutBuff *Output
	execError  string
	exitCode   int
	mutex      *sync.Mutex
}

// NewCommandWrapper is used to return a default CommandWrapper object.
func NewCommandWrapper(cmd *exec.Cmd, outBuff *Output) *CommandWrapper {
	return &CommandWrapper{
		Command:    cmd,
		StartTime:  time.Now(),
		endTime:    nil,
		StdoutBuff: outBuff,
		exitCode:   cmd.ProcessState.ExitCode(),
		mutex:      &sync.Mutex{},
	}
}

// SetEndTime sets the time when the process has finished execution.
func (cw *CommandWrapper) SetEndTime(endTime time.Time) {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()
	cw.endTime = &endTime
}

// GetEndTime gets the time when the process has finished execution.
func (cw *CommandWrapper) GetEndTime() *time.Time {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()
	return cw.endTime
}

// GetExitCode is used to access exitCode value.
func (cw *CommandWrapper) GetExitCode() int {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()
	return cw.exitCode
}

// SetExitCode is used set the exitCode value.
func (cw *CommandWrapper) SetExitCode(code int) {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()
	cw.exitCode = code
}

// GetExecError returns errors attributed to the CommandWrapper lifecycle.
func (cw *CommandWrapper) GetExecError() string {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()
	return cw.execError
}

// SetExecError appends an error attributed to the CommandWrapper lifecycle.
func (cw *CommandWrapper) SetExecError(error string) {
	cw.mutex.Lock()
	defer cw.mutex.Unlock()
	cw.execError = error
}
