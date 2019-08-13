package main

import (
  "os/exec"
  "os"
  "strconv"
  "time"
  "io"
  "bytes"
)

// handle exec calls

// TODO prune finished tasks when some max map size is reached

// maps task_id to cmd objects
var TaskIdCommandMap = make(map[string]*CommandWrapper)

// used for task_id
var Hostname, err = os.Hostname()

type CommandWrapper struct {
	Command   *exec.Cmd // underlying command
	Finished  bool // set upon process finish
  StartTime time.Time
  EndTime   time.Time
  Output    *bytes.Buffer
}

func GetProcessStatus(task_id string) StatusResponse {

  var s StatusResponse

  s.Task_id = task_id
  s.StartTime = TaskIdCommandMap[task_id].StartTime
  s.Finished = TaskIdCommandMap[task_id].Finished

  // cmd.Wait() has finished, append
  if(s.Finished) {
    s.EndTime = TaskIdCommandMap[task_id].EndTime
    s.ExitCode = TaskIdCommandMap[task_id].Command.ProcessState.ExitCode()
  }

  s.Output = TaskIdCommandMap[task_id].Output.String()
  return s
}

// start process
func RunCommand(command string) StartResponse {

  var s StartResponse

  cmd := exec.Command(command)

  var outBuf bytes.Buffer
  cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

  err := cmd.Start()

  if err != nil {
    s.Error = err.Error()
    return s
  }

  task_id := Hostname + "-" + strconv.Itoa(cmd.Process.Pid) // TODO handle if Process or Pid nil
  s.Task_id = task_id

  TaskIdCommandMap[task_id] = &CommandWrapper{cmd,false,time.Now(),time.Time{}, &outBuf}

  // async subroutine
  go func() {
      err = cmd.Wait()
      TaskIdCommandMap[task_id].Finished = true
      TaskIdCommandMap[task_id].EndTime = time.Now()
      if err != nil  {
        // TODO handle error
        s.Error = err.Error()
      }
  }()

  return s
}

// @param pid - stops process with this pid
func StopProcess(task_id string) StopResponse {

  var s StopResponse
  s.Task_id = task_id

  // check if task_id in map
  if commandWrapper, ok := TaskIdCommandMap[task_id]; ok {
    //do something here
    err := commandWrapper.Command.Process.Kill()
    if err != nil {
      s.Error = err.Error()
      return s
    } else {
      s.ExitCode = TaskIdCommandMap[task_id].Command.ProcessState.ExitCode()
      return s
    }
  }

  // task_id invalid
  s.Error = "invalid task_id"

  return s
}
