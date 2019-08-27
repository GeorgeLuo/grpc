package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GeorgeLuo/grpc/models"
	"github.com/gorilla/mux"
)

// StatusHandler returns status of running service.
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskID := vars["task_id"]

	ProcessStatusResponse, err := GetProcessStatus(taskID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorMessage{nil, err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProcessStatusResponse)
}

// StopHandler handles the logic of a call to /stop to stop a process.
func StopHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	bodyMap := make(map[string]string)
	json.Unmarshal(body, &bodyMap)

	taskID := bodyMap["task_id"]

	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorMessage{nil, "no task_id provided"})
		return
	}

	StopResponse, err := StopProcess(taskID)
	if err != nil {
		// TODO handle different error cases
		w.WriteHeader(http.StatusExpectationFailed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorMessage{&taskID, err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StopResponse)
}

// StartHandler handles the logic of a call to /start to start a process.
func StartHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	bodyMap := make(map[string]string)
	json.Unmarshal(body, &bodyMap)

	command := bodyMap["command"]
	if command == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorMessage{nil, "no command provided"})
		return
	}
	RunCommandResponse, err := RunCommand(command)

	if err != nil {
		// TODO handle different error cases, namely separate invalid task_id error code
		// though this is not terribly illogical as a response
		w.WriteHeader(http.StatusExpectationFailed)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorMessage{nil, err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RunCommandResponse)
}
