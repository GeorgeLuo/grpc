package main
import (
  "time"
)

// request and response structures

// defines object returned to status handler
type StatusResponse struct {
  Task_id string `json:"task_id"`
  Finished bool `json:"finished"`
  StartTime time.Time `json:"start_time"`
  EndTime   time.Time `json:"end_time,omitempty"`
  ExitCode  int `json:"exit_code,omitempty"`
  Output   string `json:"output"`
  Error   string `json:"error,omitempty"`
}

type StartResponse struct {
  Task_id string `json:"task_id,omitempty"`
  Error   string `json:"error,omitempty"`
}

type StopResponse struct {
  Task_id string `json:"task_id"`
  ExitCode  int `json:"exit_code,omitempty"`
  Error   string `json:"error,omitempty"`
}

type ErrorMessage struct {
  Error   string `json:"error,omitempty"`
}
