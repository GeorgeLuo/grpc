package main

import (
	"time"
)

// request and response structures

// Response structure for /status call.
type StatusResponse struct {
	Task_id   string     `json:"task_id"`
	Finished  *bool      `json:"finished,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	ExitCode  *int       `json:"exit_code,omitempty"`
	Output    []string   `json:"output,omitempty"`
	Errors    []string   `json:"errors,omitempty"`
}

// Response structure for /start call.
type StartResponse struct {
	Task_id string `json:"task_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Response structure for /stop call.
type StopResponse struct {
	Task_id  string   `json:"task_id"`
	ExitCode *int     `json:"exit_code,omitempty"`
	Error    []string `json:"errors,omitempty"`
}

// Response structure for error message.
type ErrorMessage struct {
	Error string `json:"error,omitempty"`
}
