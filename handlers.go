package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	var statusRequest models.StatusRequest
	err = json.Unmarshal(body, &statusRequest)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	taskID := statusRequest.TaskID

	ProcessStatusResponse, err := GetProcessStatus(taskID)

	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{TaskID: &taskID, Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ProcessStatusResponse)
	if err != nil {
		log.Printf("StatusHandler failed to encode with: [%s]", err.Error())
		return
	}
}

// StopHandler handles the logic of a call to /stop to stop a process.
func StopHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	var stopRequest models.StopRequest
	err = json.Unmarshal(body, &stopRequest)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	taskID := stopRequest.TaskID

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
	err = json.NewEncoder(w).Encode(StopResponse)
	if err != nil {
		log.Printf("StopHandler failed to encode with: [%s]", err.Error())
		return
	}
}

// StartHandler handles the logic of a call to /start to start a process.
func StartHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	var startRequest models.StartRequest
	err = json.Unmarshal(body, &startRequest)
	if err != nil {
		replyWithError(w, http.StatusBadRequest,
			models.ErrorMessage{Error: err.Error()})
		return
	}

	command := startRequest.Command
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
	err = json.NewEncoder(w).Encode(RunCommandResponse)
	if err != nil {
		log.Printf("StartHandler failed to encode with: [%s]", err.Error())
		return
	}
}

// helper function to return error response
func replyWithError(writer http.ResponseWriter,
	statusCode int, error models.ErrorMessage) {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(error)
	if err != nil {
		log.Printf("replyWithError failed to encode with: [%s]", err.Error())
		return
	}
}
