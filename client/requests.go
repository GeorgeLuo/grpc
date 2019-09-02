package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GeorgeLuo/grpc/models"
)

// TODO consider refactor of these 3 methods. At the moment this makes sense,
// but custom logic might be needed in the future.

// StartRequest returns a request object to start a process from /start endpoint.
func StartRequest(request models.StartRequest, host string) (*http.Request, error) {
	return buildRequest(request, startURL(host))
}

func startURL(host string) string {
	return "https://" + host + ":8443/start"
}

// StopRequest returns a request object to stop a process using /stop endpoint.
func StopRequest(request models.StopRequest, host string) (*http.Request, error) {
	return buildRequest(request, stopURL(host))
}

func stopURL(host string) string {
	return "https://" + host + ":8443/stop"
}

// StatusRequest returns a request object to retrieve status of a process
// using /status endpoint.
func StatusRequest(request models.StatusRequest, host string) (*http.Request, error) {
	return buildRequest(request, statusURL(host))
}

func statusURL(host string) string {
	return "https://" + host + ":8443/status"
}

func buildRequest(request interface{}, urlString string) (*http.Request, error) {
	byteBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("malformed body")
	}
	req, err := http.NewRequest("POST", urlString, bytes.NewReader(byteBody))
	req.Header.Set("Content-Type", "application/json")
	return req, err
}
