package models

import (
	"time"
)

// Tableble is an interface defined for structs with methods to extract
// tablewriter compliant collections.
type Tableble interface {
	Data() map[string]string
}

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

// StartRequest structure encapsulates body fields for start enpoint.
type StartRequest struct {
	Command string `json:"command"`
}

// StopRequest structure encapsulates body fields for stop enpoint.
type StopRequest struct {
	TaskID string `json:"task_id"`
}

// StatusRequest structure encapsulates body fields for status enpoint.
type StatusRequest struct {
	TaskID string `json:"task_id"`
}

// StatusBatchRequest structure encapsulates body fields for status batch
// endpoint as an array of TaskID values.
type StatusBatchRequest struct {
	TaskIDs []string `json:"task_ids"`
}

// StatusBatchResponse structure encapsulates response of the status batch
// endpoint as an array of status responses and an array of their corresponding
// errors.
type StatusBatchResponse struct {
	StatusResponses []StatusResponse `json:"status_responses"`
	Errors          []ErrorMessage   `json:"errors"`
}
