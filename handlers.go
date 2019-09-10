package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GeorgeLuo/grpc/models"
)

// StatusHandler returns status of running service.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	// TODO: make use of request struct from models.go
	bodyMap := make(map[string]string)
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	taskID := bodyMap["task_id"]

	ProcessStatusResponse, err := GetProcessStatus(taskID)

	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{TaskID: &taskID, Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProcessStatusResponse)
}

// StopHandler handles the logic of a call to /stop to stop a process.
func StopHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	// TODO: make use of request struct from models.go
	bodyMap := make(map[string]string)
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	taskID := bodyMap["task_id"]

	if taskID == "" {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: "no task_id provided"})
		return
	}

	StopResponse, err := StopProcess(taskID)
	if err != nil {
		// TODO handle different error cases
		replyWithError(w, http.StatusExpectationFailed,
			models.ErrorMessage{TaskID: &taskID, Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StopResponse)
}

// StartHandler handles the logic of a call to /start to start a process.
func StartHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	// TODO: make use of request struct from models.go
	bodyMap := make(map[string]string)
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	command := bodyMap["command"]
	if command == "" {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: "no command provided"})
		return
	}

	RunCommandResponse, err := RunCommand(command)

	if err != nil {
		// TODO handle different error cases, namely separate invalid task_id
		// error code though this is not terribly illogical as a response
		replyWithError(w, http.StatusExpectationFailed,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RunCommandResponse)
}

// helper function to return error response
func replyWithError(writer http.ResponseWriter,
	statusCode int, error models.ErrorMessage) {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(error)
}
