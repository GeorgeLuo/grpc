// handlers_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	var body = []byte(`{"task_id":"abc"}`)
	APIRouter := newAPIRouter()

	var errorMessage models.ErrorMessage
	status := getExpectedError(t, "status", body, &errorMessage, APIRouter)

	// Check the status code is what we expect.
	if status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expectedError := "invalid task_id"

	if errorMessage.Error != expectedError {
		t.Errorf("handler returned unexpected body: got %v want %v",
			errorMessage.Error, expectedError)
	}
}

// TestStartProcessBasic validates start command with status check
// after start call.
func TestStartProcessBasic(t *testing.T) {

	command := "echo 12345"
	APIRouter := newAPIRouter()

	var startResponse models.StartResponse
	status := startProcess(t, command, &startResponse, APIRouter)

	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	time.Sleep(1 * time.Second) // now check status

	var statusResponse models.StatusResponse
	status = processStatus(t, startResponse.TaskID, &statusResponse, APIRouter)

	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedStatusOutput := []string{"12345"}

	if !reflect.DeepEqual(statusResponse.Output, expectedStatusOutput) {
		t.Errorf("command output returned unexpected value: got %v want %v",
			statusResponse.Output, expectedStatusOutput)
	}

	if *statusResponse.ExitCode != 0 {
		t.Errorf("command exit code returned unexpected value: got %d want %d",
			statusResponse.ExitCode, 0)
	}
}

// TestStopProcessBasic validates stop command with status check after
// start then stop call.
func TestStopProcessBasic(t *testing.T) {

	command := "sleep 5"
	APIRouter := newAPIRouter()

	var startResponse models.StartResponse
	status := startProcess(t, command, &startResponse, APIRouter)

	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	time.Sleep(1 * time.Second) // now stop process

	var stopResponse models.StopResponse
	status = stopProcess(t, startResponse.TaskID, &stopResponse, APIRouter)

	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if *stopResponse.ExitCode != -1 {
		t.Errorf("command exit code returned unexpected value: got %d want %d",
			stopResponse.ExitCode, -1)
	}
}

// Helper method to send start command.
func startProcess(t *testing.T, command string,
	startResponse *models.StartResponse, r *mux.Router) int {

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(models.StartRequest{Command: command})

	startRequest, err := http.NewRequest("POST", "/start",
		bytes.NewBuffer(body.Bytes()))

	if err != nil {
		t.Fatal(err)
	}
	startRequest.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, startRequest)

	json.Unmarshal([]byte(rr.Body.String()), startResponse)
	return rr.Code
}

// Helper method to send stop request.
func stopProcess(t *testing.T, taskID string,
	stopResponse *models.StopResponse, r *mux.Router) int {

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(models.StopRequest{TaskID: taskID})

	stopRequest, err := http.NewRequest("POST", "/stop",
		bytes.NewBuffer(body.Bytes()))

	if err != nil {
		t.Fatal(err)
	}
	stopRequest.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, stopRequest)

	json.Unmarshal([]byte(rr.Body.String()), stopResponse)
	return rr.Code
}

// Helper method to send stop request.
func processStatus(t *testing.T, taskID string,
	statusResponse *models.StatusResponse, r *mux.Router) int {

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(models.StatusRequest{TaskID: taskID})

	statusRequest, err := http.NewRequest("POST", "/status",
		bytes.NewBuffer(body.Bytes()))

	if err != nil {
		t.Fatal(err)
	}
	statusRequest.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, statusRequest)

	json.Unmarshal([]byte(rr.Body.String()), statusResponse)
	return rr.Code
}

// Helper method to make calls that return a generic error.
func getExpectedError(t *testing.T, endpoint string, body []byte,
	errorResponse *models.ErrorMessage, r *mux.Router) int {

	request, err := http.NewRequest("POST", "/"+endpoint, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, request)

	json.Unmarshal([]byte(rr.Body.String()), errorResponse)
	return rr.Code
}

// Helper method to provide a global context router.
func newAPIRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/start", StartHandler).
		Methods("POST")
	r.HandleFunc("/stop", StopHandler).
		Methods("POST")
	r.HandleFunc("/status", StatusHandler).
		Methods("POST")
	return r
}