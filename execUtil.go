package main

import (
  "os/exec"
  "os"
  "log"
  "strconv"
  "time"
  "io"
  "bytes"
)

// handle exec calls

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

  log.Printf("GetProcessStatus task_id=%s, pid=%d, ExitCode=%d, Finished=%v",
    task_id, TaskIdCommandMap[task_id].Command.Process.Pid,
    TaskIdCommandMap[task_id].Command.ProcessState.ExitCode(), TaskIdCommandMap[task_id].Finished)

  // return TaskIdCommandMap[task_id].Command.ProcessState.ExitCode()
  return s
}

// start process
func RunCommand(command string) string {
  log.Printf("command=%s", command)

  cmd := exec.Command(command)

  var outBuf bytes.Buffer
  cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

  err := cmd.Start()
  if err != nil {
     log.Fatal(err)
  }

  task_id := Hostname + "-" + strconv.Itoa(cmd.Process.Pid) // TODO handle if Process or Pid nil
  TaskIdCommandMap[task_id] = &CommandWrapper{cmd,false,time.Now(),time.Time{}, &outBuf}

  // async subroutine
  go func() {
      log.Printf("pre-wait ExitCode=%v", cmd.ProcessState.ExitCode()) // ExitCode can be -1 when stilling running
      err = cmd.Wait()
      TaskIdCommandMap[task_id].Finished = true
      TaskIdCommandMap[task_id].EndTime = time.Now()
      if err != nil  {
        // TODO handle error
        log.Printf("Command finished with error: %v", err)
      } else {
        log.Printf("post-wait ExitCode=%v", cmd.ProcessState.ExitCode())
      }
  }()

  log.Printf("RunCommand task_id=%s, pid=%d", task_id, TaskIdCommandMap[task_id].Command.Process.Pid)
  return task_id
}

// @param pid - stops process with this pid
func StopProcess(task_id string) int {

  // check if already finished

  err := TaskIdCommandMap[task_id].Command.Process.Kill()
  if err != nil {
     log.Fatal(err)
  }
  log.Printf("StopProcess task_id=%s, pid=%d", task_id, TaskIdCommandMap[task_id].Command.Process.Pid)
  return TaskIdCommandMap[task_id].Command.ProcessState.ExitCode()
}
