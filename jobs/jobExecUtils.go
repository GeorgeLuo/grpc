package jobs

import (
	"errors"

	"github.com/GeorgeLuo/grpc/core"
	"github.com/GeorgeLuo/grpc/models"
	"github.com/GeorgeLuo/grpc/utils"
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
		startResponse, err := core.RunCommand(request.Command, request.Alias)
		if err != nil {
			bailErr := stopCommands(successTaskIDs)
			if len(bailErr) > 0 {
				err = utils.AppendStringToError(err, bailErr...)
			}
			return nil, err
		}

		successTaskIDs = append(successTaskIDs, startResponse.TaskID)
		jobStartResponse.StartResponses = append(jobStartResponse.StartResponses,
			*startResponse)
	}

	core.GlobalAliasMap.Put(alias, successTaskIDs...)
	return &jobStartResponse, nil
}

// GetJobStatusByAlias retrieves status of a running job
func GetJobStatusByAlias(alias string) (*models.JobStatusResponse, error) {

	taskIDs, ok := core.GlobalAliasMap.Get(alias)
	if !ok {
		return nil, errors.New("alias not mapped")
	}

	var jobStatusResponse models.JobStatusResponse
	var statusErrors []string

	for _, taskID := range taskIDs {
		statusResponse, err := core.GetProcessStatus(taskID)
		if err != nil {
			statusErrors = append(statusErrors, err.Error())
		} else {
			jobStatusResponse.StatusResponses =
				append(jobStatusResponse.StatusResponses, *statusResponse)
		}
	}

	jobStatusResponse.Errors = statusErrors
	return &jobStatusResponse, nil
}

// StopJobByAlias stops all tasks linked to the job alias.
func StopJobByAlias(alias string) (*models.JobStopResponse, error) {

	taskIDs, ok := core.GlobalAliasMap.Get(alias)
	if !ok {
		return nil, errors.New("alias not mapped")
	}

	var jobStopResponse models.JobStopResponse

	jobStopResponse.Errors = stopCommands(taskIDs)
	return &jobStopResponse, nil
}

// terminates processes with taskIDs provided
func stopCommands(taskIDs []string) []string {
	var errors []string

	for _, taskID := range taskIDs {
		_, err := core.StopProcess(taskID)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}
	return errors
}
