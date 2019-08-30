// handlers_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/GeorgeLuo/grpc/models"
	"github.com/gorilla/mux"
)

// TODO test cases
// start long process, get statusResponse, wait until finished, get statusResponse
// start non-existant process
// start process without permission
// start process that returns error code, get statusResponse
// start long process, stop process, get statusResponse
// stop with invalid task_id
// stop already finished process

// TestStatusWithInvalidTaskID tests case of invalid status
func TestStatusWithInvalidTaskID(t *testing.T) {

	var data = []byte(`{"task_id":"abc"}`)

	req, err := http.NewRequest("POST", "/status", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler).Methods("POST")

	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"invalid task_id"}`

	if strings.TrimSuffix(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestStartProcessBasic validates start command with status check after start call.
func TestStartProcessBasic(t *testing.T) {

	var data = []byte(`{"command":"echo 12345"}`)

	startRequest, err := http.NewRequest("POST", "/start", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	startRequest.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler).
		Methods("POST")
	r.HandleFunc("/start", StartHandler).
		Methods("POST")

	r.ServeHTTP(rr, startRequest)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	startResponse := &models.StartResponse{}
	json.Unmarshal([]byte(rr.Body.String()), startResponse)
	taskID := startResponse.TaskID

	time.Sleep(1 * time.Second) // now check status

	data = []byte(`{"task_id":"` + taskID + `"}`)

	req, err := http.NewRequest("POST", "/status", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	statusResponse := &models.StatusResponse{}
	json.Unmarshal([]byte(rr.Body.String()), statusResponse)

	expectedStatusOutput := []string{"12345"}

	if !reflect.DeepEqual(statusResponse.Output, expectedStatusOutput) {
		t.Errorf("command output returned unexpected value: got %v want %v",
			statusResponse.Output, expectedStatusOutput)
	}

	if *statusResponse.ExitCode != 0 {
		t.Errorf("command exit code returned unexpected value: got %d want %d",
			statusResponse.ExitCode, 0)
	}

	if !*statusResponse.Finished {
		t.Errorf("command finished returned unexpected value: got %t want %t",
			*statusResponse.Finished, true)
	}
}

// TestStopProcessBasic validates stop command with status check after start then stop call.
func TestStopProcessBasic(t *testing.T) {

	var data = []byte(`{"command":"sleep 5"}`)

	startRequest, err := http.NewRequest("POST", "/start", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	startRequest.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/start", StartHandler).
		Methods("POST")
	r.HandleFunc("/stop", StopHandler).
		Methods("POST")

	r.ServeHTTP(rr, startRequest)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	startResponse := &models.StartResponse{}
	json.Unmarshal([]byte(rr.Body.String()), startResponse)
	taskID := startResponse.TaskID

	time.Sleep(1 * time.Second) // now stop process

	data = []byte(`{"task_id":"` + taskID + `"}`)

	stopRequest, err := http.NewRequest("POST", "/stop", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, stopRequest)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	stopResponse := &models.StopResponse{}
	json.Unmarshal([]byte(rr.Body.String()), stopResponse)

	if *stopResponse.ExitCode != -1 {
		t.Errorf("command exit code returned unexpected value: got %d want %d",
			stopResponse.ExitCode, -1)
	}
}
