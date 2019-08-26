package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// handle exec calls

// TODO prune finished tasks when some max map size is reached

var taskIDCommandMap SyncMap
var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	taskIDCommandMap = NewMap()
}

// GetProcessStatus is called to retrieve the details of a processes by task_id.
func GetProcessStatus(taskID string) (*StatusResponse, error) {

	var statusResponse StatusResponse
	statusResponse.TaskID = taskID

	var command *CommandWrapper
	var ok bool

	// validate task_id
	if command, ok = taskIDCommandMap.Get(taskID); !ok {
		return nil, errors.New("invalid task_id")
	}

	// TODO refactor this into un/finished process
	statusResponse.TaskID = taskID
	statusResponse.StartTime = new(time.Time)
	*statusResponse.StartTime = command.StartTime
	statusResponse.ExecError = command.GetExecError()
	statusResponse.Finished = new(bool)
	if command.GetEndTime() != nil {
		*statusResponse.Finished = true
		statusResponse.EndTime = new(time.Time)
		statusResponse.EndTime = command.GetEndTime()
		statusResponse.ExitCode = new(int)
		*statusResponse.ExitCode = command.GetExitCode()
	}

	statusResponse.Output = command.StdoutBuff.Lines()

	return &statusResponse, nil
}

// RunCommand starts a process from command argument.
func RunCommand(command string) (*StartResponse, error) {

	var startResponse StartResponse
	splitCommand := strings.Split(command, " ")

	cmd := exec.Command(splitCommand[0], splitCommand[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	outBuf := NewOutput()
	cmd.Stdout = io.MultiWriter(os.Stdout, outBuf)

	err := cmd.Start()

	if err != nil {
		return nil, err
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)

	if err != nil {
		return nil, err
	}

	taskID := hostname + "-" + strconv.Itoa(pgid) // TODO handle if Process or Pid nil
	startResponse.TaskID = taskID

	taskIDCommandMap.Put(taskID, NewCommandWrapper(cmd, outBuf))

	go func() {
		// TODO add append error to CommandWrapper, impl accessors and setters
		waitErr := cmd.Wait()
		cw, _ := taskIDCommandMap.Get(taskID)
		if waitErr != nil {
			cw.SetExecError(waitErr.Error())
		}
		cw.SetExitCode(cmd.ProcessState.ExitCode())
		cw.SetEndTime(time.Now())
	}()

	return &startResponse, nil
}

func intPtr(value int) *int {
	return &value
}

// StopProcess is called to end a previously started process.
func StopProcess(taskID string) (*StopResponse, error) {

	var stopResponse StopResponse

	stopResponse.TaskID = taskID

	var command *CommandWrapper
	var ok bool

	// if task_id invalid, return error.
	if command, ok = taskIDCommandMap.Get(taskID); !ok {
		return &stopResponse, errors.New("invalid task_id")
	}

	if command.GetEndTime() != nil {
		stopResponse.ExitCode = intPtr(command.GetExitCode())
		return &stopResponse, errors.New("process already finished")
	}

	// send stop signal first
	// TODO verify flat processes work with gpid as well

	pgid, _ := syscall.Getpgid(command.Command.Process.Pid)
	err := syscall.Kill(-pgid, syscall.SIGQUIT)

	if err != nil && err != syscall.EPERM {

		err = syscall.Kill(-pgid, syscall.SIGKILL)
		if err != nil {
			return &stopResponse, err
		}
	}

	stopResponse.ExitCode = intPtr(command.GetExitCode())

	return &stopResponse, err
}
