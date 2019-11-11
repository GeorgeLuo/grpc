package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GeorgeLuo/grpc/core"
	"github.com/GeorgeLuo/grpc/models"
)

// these handlers are responsible for alias mapped operations
// the aliases from the client and are compared on job start for conflicted
// names. From then on the jobs can be accessed much as the coreHandlers,
// with the alias instead of TaskID.

// JobStartHandler is used to start a process with an alias
func JobStartHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		core.ReplyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	var jobStartRequest models.JobStartRequest
	err = json.Unmarshal(body, &jobStartRequest)
	if err != nil {
		core.ReplyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	startRequests := jobStartRequest.StartRequests
	for i, request := range startRequests {
		if request.Command == "" {
			errMsg := fmt.Sprintf("command missing at index: %d", i)
			core.ReplyWithError(w, http.StatusBadRequest,
				models.ErrorMessage{Error: errMsg})
			return
		}
	}

	runJobResponse, err := RunJob(jobStartRequest.StartRequests,
		jobStartRequest.Alias)

	if err != nil {
		// TODO handle different error cases, namely separate invalid task_id
		// error code though this is not terribly illogical as a response
		core.ReplyWithError(w, http.StatusExpectationFailed,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(runJobResponse)
	if err != nil {
		log.Printf("StartHandler failed to encode with: [%s]", err.Error())
		return
	}
}

// JobStatusHandler returns status of running job.
func JobStatusHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		core.ReplyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	var jobStatusRequest models.JobStatusRequest
	err = json.Unmarshal(body, &jobStatusRequest)
	if err != nil {
		core.ReplyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	alias := jobStatusRequest.Alias

	if alias == "" {
		core.ReplyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: "no alias provided"})
		return
	}

	var jobStatusResponse *models.JobStatusResponse
	jobStatusResponse, err = GetJobStatusByAlias(jobStatusRequest.Alias)

	if err != nil {
		core.ReplyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	if len(jobStatusResponse.Errors) > 0 {
		w.WriteHeader(http.StatusExpectationFailed)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jobStatusResponse)
	if err != nil {
		log.Printf("StatusHandler failed to encode with: [%s]", err.Error())
		return
	}
}
