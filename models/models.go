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

// JOB MODEL DEFINITIONS BEGIN

// JobStartResponse structure encapsulates body fields for job start response,
// a wrapper for all the associated tasks.
type JobStartResponse struct {
	StartResponses []StartResponse `json:"tasks,omitempty"`
	Alias          string          `json:"alias,omitempty"`
}

// JobStatusResponse structure encapsulates body fields for job status,
// a wrapper for all the associated tasks.
type JobStatusResponse struct {
	StatusResponses []StatusResponse `json:"tasks,omitempty"`
	Alias           string           `json:"alias,omitempty"`
	Errors          []error          `json:"errors,omitempty"`
}

// TODO: The []StartRequest can be reduced to []string, keeping the object
// array for now in case StartRequest will contain more metadata
// (Scheduling parameters).

// JobStartRequest structure wraps multiple start requests with an associated
// alias.
type JobStartRequest struct {
	StartRequests []StartRequest `json:"tasks,omitempty"`
	Alias         string         `json:"alias,omitempty"`
}

// JobStatusRequest structure wraps multiple status requests with an associated
// alias.
type JobStatusRequest struct {
	Alias string `json:"alias,omitempty"`
}
