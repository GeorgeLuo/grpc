package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// StatusHandler returns status of running service.
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	task_id := vars["task_id"]

	res := GetProcessStatus(task_id)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// StopHandler handles the logic of a call to /stop to stop a process.
func StopHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	body_map := make(map[string]string)
	json.Unmarshal(body, &body_map)

	task_id := body_map["task_id"]

	if task_id == "" {
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorMessage{"no task_id provided"})
		return
	}

	StopResponse := StopProcess(task_id)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StopResponse)

}

// StartHandler handles the logic of a call to /start to start a process.
func StartHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	body_map := make(map[string]string)
	json.Unmarshal(body, &body_map)

	command := body_map["command"]
	if command == "" {
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorMessage{"no command provided"})
		return
	}
	res := RunCommand(command)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
