package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// path param task_id : identifier of task query
// returns status of running service
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	task_id := vars["task_id"]

	res := AsyncGetProcessStatus(task_id)
	// statusResponse := AsyncGetProcessStatus(task_id)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

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

func AsyncStartHandler(w http.ResponseWriter, r *http.Request) {

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
	res := AsyncRunCommand(command)

	// startResponse := RunCommand(command)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
