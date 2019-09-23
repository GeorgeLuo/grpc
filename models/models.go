package models

import (
	"strconv"
	"time"
)

// TODO: handle batch request and responses

// Renderable is an interface to support operation to display contents in a
// tablewriter table.
type Renderable interface {
	Headers() []string
	Rows() [][]string
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

// Headers returns the headers to populate a table of status response fields.
func (r *StatusResponse) Headers() []string {
	return []string{"task_id", "start_time", "end_time",
		"exit_code", "exec_error"}
}

// Rows produces a row of data for the data returned by status response.
func (r *StatusResponse) Rows() [][]string {
	return [][]string{
		{r.TaskID, r.StartTime.String(), r.EndTime.String(),
			strconv.Itoa(*r.ExitCode), r.ExecError},
	}
}

// StartResponse structure to capture return for /start call.
type StartResponse struct {
	TaskID string `json:"task_id,omitempty"`
}

// Headers returns the headers to display start response data.
func (r *StartResponse) Headers() []string {
	return []string{"task_id"}
}

// Rows produces a row of data for task_id.
func (r *StartResponse) Rows() [][]string {
	return [][]string{
		{r.TaskID},
	}
}

// StopResponse is a structure to capture return for /stop call.
type StopResponse struct {
	TaskID   string `json:"task_id"`
	ExitCode *int   `json:"exit_code,omitempty"`
}

// Headers returns the headers to populate a table of stop response fields.
func (r *StopResponse) Headers() []string {
	return []string{"task_id", "exit_code"}
}

// Rows produces a row of data for the data returned by stop response.
func (r *StopResponse) Rows() [][]string {
	return [][]string{
		{r.TaskID, strconv.Itoa(*r.ExitCode)},
	}
}

// ErrorMessage is structure for error message.
type ErrorMessage struct {
	TaskID *string `json:"task_id,omitempty"`
	Error  string  `json:"error,omitempty"`
}

// StartRequest structure encapsulates body fields for start enpoint.
type StartRequest struct {
	Command string `json:"command,omitempty"`
}

// StopRequest structure encapsulates body fields for stop enpoint.
type StopRequest struct {
	TaskID string `json:"task_id,omitempty"`
}

// StatusRequest structure encapsulates body fields for status enpoint.
type StatusRequest struct {
	TaskID string `json:"task_id,omitempty"`
}
