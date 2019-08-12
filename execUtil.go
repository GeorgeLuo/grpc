package main

import (
  "os/exec"
  "os"
  "log"
  "strconv"
)

// handle exec calls


// maps task_id to cmd objects
var TaskIdCmdMap = make(map[string]*exec.Cmd)

// TODO define Command class

func GetProcessStatus(task_id string) int {
  log.Printf("GetProcessStatus task_id=%s, pid=%d, ProcessState=%d", task_id, TaskIdCmdMap[task_id].Process.Pid, TaskIdCmdMap[task_id].ProcessState)
  return TaskIdCmdMap[task_id].ProcessState.ExitCode()
}

// start process
func RunCommand(command string) string {
  log.Printf("command=%s", command)
  cmd := exec.Command(command)
  cmd.Stdout = os.Stdout
  err := cmd.Start()
  if err != nil {
     log.Fatal(err)
  }

  go func() {
      log.Printf("pre-wait ProcessState=%v", cmd.ProcessState.ExitCode())
      err = cmd.Wait()
      log.Printf("Command finished with error: %v", err)

      // TODO mark cmd as finished in subroutine
  }()
  hostname, err := os.Hostname()

  // for now make this for unique task_id
  task_id := hostname + "-" + strconv.Itoa(cmd.Process.Pid)
  TaskIdCmdMap[task_id] = cmd

  log.Printf("RunCommand task_id=%s, pid=%d", task_id, TaskIdCmdMap[task_id].Process.Pid)

  return task_id
}

// @param pid - stops process with this pid
func StopProcess(task_id string) int {

  // check if already finished

  err := TaskIdCmdMap[task_id].Process.Kill()
  if err != nil {
     log.Fatal(err)
  }
  log.Printf("StopProcess task_id=%s, pid=%d", task_id, TaskIdCmdMap[task_id].Process.Pid)
  return TaskIdCmdMap[task_id].ProcessState.ExitCode()
}
