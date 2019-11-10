package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/GeorgeLuo/grpc/models"
)

// handle exec calls

// TODO prune finished tasks when some max map size is reached

var taskIDCommandMap SyncMap
var hostname string

// GlobalAliasMap is used to retrieve task_ids started under an alias
var GlobalAliasMap AliasMap

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	taskIDCommandMap = NewMap()
	GlobalAliasMap = NewAliasMap()
}

// GetProcessStatusByAlias retrieves status using alias mapped to task_id(s)
func GetProcessStatusByAlias(alias string) (*models.StatusResponse, error) {
	if taskID, ok := GlobalAliasMap.Get(alias); ok {
		return GetProcessStatus(taskID[0])
	}

	return nil, errors.New("alias not mapped")
}

// GetProcessStatus retrieves the status of the process specified with taskID.
func GetProcessStatus(taskID string) (*models.StatusResponse, error) {

	var statusResponse models.StatusResponse

	var command *CommandWrapper
	var ok bool

	// validate task_id
	if command, ok = taskIDCommandMap.Get(taskID); !ok {
		return nil, errors.New("no process mapped to task_id")
	}

	// TODO refactor this into un/finished process
	statusResponse.TaskID = taskID
	statusResponse.StartTime = new(time.Time)
	*statusResponse.StartTime = command.StartTime
	statusResponse.ExecError = command.GetExecError()
	if command.GetEndTime() != nil {
		statusResponse.EndTime = new(time.Time)
		statusResponse.EndTime = command.GetEndTime()
		statusResponse.ExitCode = new(int)
		*statusResponse.ExitCode = command.GetExitCode()
	}

	statusResponse.Output = command.StdoutBuff.GetOutput()

	return &statusResponse, nil
}

// RunCommand starts a process from command argument.
func RunCommand(command string, alias string) (*models.StartResponse, error) {

	var startResponse models.StartResponse
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

	// TODO: perhaps handle duplicate taskID with error from Put
	taskIDCommandMap.Put(taskID, NewCommandWrapper(cmd, outBuf))

	if alias != "" {
		err = GlobalAliasMap.Put(alias, taskID)
		if err != nil {
			return nil, fmt.Errorf("no process started, alias map err: %s",
				err.Error())
		}
	}

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

// StopProcessByAlias stops process using alias mapped to task_id(s)
func StopProcessByAlias(alias string) (*models.StopResponse, error) {
	if taskID, ok := GlobalAliasMap.Get(alias); ok {
		return StopProcess(taskID[0])
	}

	return nil, errors.New("alias not mapped")
}

// StopProcess ends a previously started process.
func StopProcess(taskID string) (*models.StopResponse, error) {

	var stopResponse models.StopResponse

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
