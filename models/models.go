package models

import (
	"time"
)

// TODO: handle batch request and responses
// TODO: make renderable wrapper around the responses instead of making
// responses renderables.

// request and response structures

// StatusResponse is a structure to capture return for /status call.
type StatusResponse struct {
	TaskID string `json:"task_id"`
	// StartTime's value is populated when the process is first called.
	StartTime *time.Time `json:"start_time,omitempty"`
	// EndTime is populated when the process is marked as Finished in ISO 8601 format.
	EndTime  *time.Time `json:"end_time,omitempty"`
	ExitCode *int       `json:"exit_code,omitempty"`
	// Output is the contents of the output buffer of the process up to the moment.
	Output []string `json:"output,omitempty"`
	// ExecError is the output taken from syscall's error response.
	ExecError string `json:"execError,omitempty"`
}

// StartResponse structure to capture return for /start call.
type StartResponse struct {
	TaskID string `json:"task_id,omitempty"`
	Alias  string `json:"alias,omitempty"`
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

// StartRequest structure encapsulates body fields for start endpoint.
type StartRequest struct {
	Command string `json:"command,omitempty"`
	Alias   string `json:"alias,omitempty"`
}

// StopRequest structure encapsulates body fields for stop endpoint.
type StopRequest struct {
	TaskID string `json:"task_id,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

// StatusRequest structure encapsulates body fields for status endpoint.
type StatusRequest struct {
	TaskID string `json:"task_id,omitempty"`
	Alias  string `json:"alias,omitempty"`
}
