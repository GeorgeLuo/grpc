package main

import (
	"errors"

	"github.com/GeorgeLuo/grpc/models"
)

// the GlobalAliasMap referenced comes from the coreExecUtil definition
// the implication is aliases can be defined to map only one of a task or a job

// RunJob starts a series of commands under the alias name, and maintains the
// alias mapping to the individual command alias underneath.
func RunJob(startRequests []models.StartRequest,
	alias string) (*models.JobStartResponse, error) {

	var jobStartResponse models.JobStartResponse
	var successTaskIDs []string

	for _, request := range startRequests {
		startResponse, err := RunCommand(request.Command, request.Alias)
		if err != nil {
			bailErr := bailCommands(successTaskIDs)
			if bailErr != nil {
				err = AppendError(err, bailErr)
			}
			return nil, err
		}

		successTaskIDs = append(successTaskIDs, startResponse.TaskID)
		jobStartResponse.StartResponses = append(jobStartResponse.StartResponses,
			*startResponse)
	}

	GlobalAliasMap.Put(alias, successTaskIDs...)
	return &jobStartResponse, nil
}

// GetJobStatusByAlias retrieves status of a running job
func GetJobStatusByAlias(alias string) (*models.JobStatusResponse, error) {

	taskIDs, ok := GlobalAliasMap.Get(alias)
	if !ok {
		return nil, errors.New("alias not mapped")
	}

	var jobStatusResponse models.JobStatusResponse
	var statusErrors []error

	for _, taskID := range taskIDs {
		statusResponse, err := GetProcessStatus(taskID)
		if err != nil {
			statusErrors = append(statusErrors, err)
		} else {
			jobStatusResponse.StatusResponses =
				append(jobStatusResponse.StatusResponses, *statusResponse)
		}
	}

	jobStatusResponse.Errors = statusErrors
	return &jobStatusResponse, nil
}

// terminates processes with taskIDs provided
func bailCommands(taskIDs []string) error {
	for _, taskID := range taskIDs {
		_, err := StopProcess(taskID)
		if err != nil {
			return err
		}
	}
	return nil
}
