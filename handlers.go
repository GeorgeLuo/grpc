
package main
import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"io/ioutil"
)

// TODO add message fields

type StartResponse struct {
    Task_id string `json:"task_id"`
}

type StopResponse struct {
    Task_id string `json:"task_id"`
		ExitCode int `json:"exit_code"`
}

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
	log.Printf("body=%s", body)

	body_map := make(map[string]string)
	json.Unmarshal(body, &body_map)

	task_id := body_map["task_id"]
	log.Printf("task_id=%s", task_id)

	// TODO define and handle errors
	output := StopProcess(task_id)

	log.Printf("task_id=%s, status=%s", task_id, output)

	responseBody := StartResponse{task_id}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)

}

// return pid
func StartHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	log.Printf("body=%s", body)

	body_map := make(map[string]string)
	json.Unmarshal(body, &body_map)

	command := body_map["command"]

	// TODO define and handle errors
	task_id := RunCommand(command)

	log.Printf("command=%s, task_id=%s", command, task_id)

	responseBody := StartResponse{task_id}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody)
}
