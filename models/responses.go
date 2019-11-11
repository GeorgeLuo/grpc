package models

import "time"

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

// JobStartResponse structure encapsulates body fields for job start response,
// a wrapper for all the associated tasks.
type JobStartResponse struct {
	Alias          string          `json:"alias,omitempty"`
	StartResponses []StartResponse `json:"tasks,omitempty"`
}

// JobStatusResponse structure encapsulates body fields for job status,
// a wrapper for all the associated tasks.
type JobStatusResponse struct {
	Alias           string           `json:"alias,omitempty"`
	StatusResponses []StatusResponse `json:"tasks,omitempty"`
	Errors          []string         `json:"errors,omitempty"`
}

// JobStopResponse structure encapsulates body fields for a job stop response,
// a wrapper for all the associated task stop output.
type JobStopResponse struct {
	Alias         string         `json:"alias,omitempty"`
	StopResponses []StopResponse `json:"tasks,omitempty"`
	Errors        []string       `json:"errors,omitempty"`
}
