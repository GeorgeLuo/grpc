
package main
import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
)

// path param task_id : identifier of task query
// returns status of running service
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	task_id := vars["task_id"]

	statusResponse := GetProcessStatus(task_id)
	// log.Printf("task_id=%s, status=%d", task_id, status)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse)
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

	// TODO define and handle errors
	startResponse := RunCommand(command)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(startResponse)
}
