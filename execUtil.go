package main

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// handle exec calls

// TODO prune finished tasks when some max map size is reached

var taskIdCommandMap SyncMap
var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	taskIdCommandMap = NewMap()
}

// GetProcessStatus is called to retrieve the details of a processes by task_id.
func GetProcessStatus(task_id string) StatusResponse {

	var s StatusResponse
	s.Task_id = task_id

	var command *CommandWrapper
	var ok bool

	// validate task_id
	if command, ok = taskIdCommandMap.Get(task_id); !ok {
		s.Errors = []string{"invalid task_id"}
		return s
	}

	// TODO refactor this into un/finished process
	s.Task_id = task_id
	s.StartTime = new(time.Time)
	*s.StartTime = command.StartTime
	s.Finished = new(bool)
	*s.Finished = command.GetFinished()
	s.Errors = command.GetErrors()

	// cmd.Wait() has finished, append
	if command.GetFinished() {
		s.EndTime = new(time.Time)
		*s.EndTime = command.EndTime
		s.ExitCode = new(int)
		*s.ExitCode = command.GetExitCode()
	}

	s.Output = command.StdoutBuff.Lines()

	return s
}

// RunCommand starts a process from command argument.
func RunCommand(command string) StartResponse {
	var s StartResponse

	cmd := exec.Command(command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	outBuf := NewOutput()
	cmd.Stdout = io.MultiWriter(os.Stdout, outBuf)

	err := cmd.Start()

	if err != nil {
		s.Error = err.Error()
		return s
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)

	if err != nil {
		s.Error = err.Error()
		return s
	}

	task_id := hostname + "-" + strconv.Itoa(pgid) // TODO handle if Process or Pid nil
	s.Task_id = task_id

	taskIdCommandMap.Put(task_id, NewCommandWrapper(cmd, outBuf))

	// async subroutine
	go func() {
		// TODO add append error to CommandWrapper, impl accessors and setters
		err = cmd.Wait()
		cw, _ := taskIdCommandMap.Get(task_id)
		if err != nil {
			cw.AppendError(err.Error())
		}
		cw.SetExitCode(cmd.ProcessState.ExitCode())
		cw.SetFinished(true)
	}()

	return s
}

// StopProcess is called to end a previously started process.
func StopProcess(task_id string) StopResponse {

	var s StopResponse
	s.Task_id = task_id

	var command *CommandWrapper
	var ok bool

	// if task_id invalid, return error.
	if command, ok = taskIdCommandMap.Get(task_id); !ok {
		s.Error = []string{"invalid task_id"}
		return s
	}

	s.ExitCode = new(int)

	if command.GetFinished() {
		s.Error = []string{"process already finished."}
		*s.ExitCode = command.GetExitCode()
		return s
	}

	// send stop signal first
	// TODO verify flat processes work with gpid as well

	pgid, err := syscall.Getpgid(command.Command.Process.Pid)
	err = syscall.Kill(-pgid, syscall.SIGINT)

	if err != nil {

		errs := []string{err.Error()}

		err = syscall.Kill(-pgid, syscall.SIGKILL)
		if err != nil {
			errs = append(errs, err.Error())
			s.Error = errs
			return s
		}

		s.Error = errs
	}

	*s.ExitCode = command.GetExitCode()

	return s
}
