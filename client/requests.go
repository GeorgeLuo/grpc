package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// TODO consider refactor of these 3 methods. At the moment this makes sense,
// but custom logic might be needed in the future.

// StartRequest forms a request object for http to consume.
func StartRequest(body interface{}, host string) (*http.Request, error) {
	urlString := "https://" + host + ":8443/start"
	byteBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("malformed body")
	}
	request, err := http.NewRequest("POST", urlString, bytes.NewReader(byteBody))
	request.Header.Set("Content-Type", "application/json")
	return request, err
}

// StopRequest forms a request object for http to consume.
func StopRequest(body interface{}, host string) (*http.Request, error) {
	urlString := "https://" + host + ":8443/stop"
	byteBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("malformed body")
	}
	request, err := http.NewRequest("POST", urlString, bytes.NewReader(byteBody))
	request.Header.Set("Content-Type", "application/json")
	return request, err
}

// StatusRequest forms a request object for http to consume.
func StatusRequest(body interface{}, host string) (*http.Request, error) {
	urlString := "https://" + host + ":8443/status"
	byteBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("malformed body")
	}
	request, err := http.NewRequest("POST", urlString, bytes.NewReader(byteBody))
	request.Header.Set("Content-Type", "application/json")
	return request, err
}
