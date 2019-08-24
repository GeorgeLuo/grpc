package main

import (
	"time"
)

// request and response structures

// StatusResponse is a structure to capture return for /status call.
// Finished indicates the completion of a process (through termination or end of process with exit code).
// StartTime's value is populated when the process is first called.
// EndTime is populated when the process is marked as Finished in ISO 8601 format.
// Output is the contents of the output buffer of the process up to the moment.
// ExecError is the output taken from syscall's error response.
type StatusResponse struct {
	TaskID    string     `json:"task_id"`
	Finished  *bool      `json:"finished,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	ExitCode  *int       `json:"exit_code,omitempty"`
	Output    []string   `json:"output,omitempty"`
	ExecError *string    `json:"execError,omitempty"`
}

// StartResponse structure to capture return for /start call.
type StartResponse struct {
	TaskID string `json:"task_id,omitempty"`
}

// StopResponse is a structure to capture return for /stop call.
type StopResponse struct {
	TaskID   string `json:"task_id"`
	ExitCode *int   `json:"exit_code,omitempty"`
}

// ErrorMessage is structure for error message.
type ErrorMessage struct {
	TaskID *string `json:"task_id,omitempty"`
	Error  string  `json:"error,omitempty"`
}
