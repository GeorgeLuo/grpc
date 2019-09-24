package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GeorgeLuo/grpc/models"
)

// Start is a client handler that builds a start request and returns a
// StartResponse.
func Start(request models.StartRequest,
	host string, permission Permission) (*models.StartResponse, error) {

	httpRequest, err := StartRequest(request, host)
	if err != nil {
		return nil,
			fmt.Errorf("error forming start request:\n %s", err.Error())
	}

	responseBody, err := send(httpRequest, permission)
	if err != nil {
		return nil,
			fmt.Errorf("error sending start request:\n %s", err.Error())
	}

	var response *models.StartResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil,
			fmt.Errorf("error on unmarshal start response:\n %s", err.Error())
	}

	return response, nil
}

// Stop is a client handler that builds a stop request and returns a
// StopResponse.
func Stop(request models.StopRequest,
	host string, permission Permission) (*models.StopResponse, error) {

	httpRequest, err := StopRequest(request, host)
	if err != nil {
		return nil,
			fmt.Errorf("error forming stop request:\n %s", err.Error())
	}

	responseBody, err := send(httpRequest, permission)
	if err != nil {
		return nil,
			fmt.Errorf("error sending stop request:\n %s", err.Error())
	}

	var response *models.StopResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil,
			fmt.Errorf("error on unmarshal stop response:\n %s", err.Error())
	}

	return response, nil
}

// Status is a client handler that builds a stop request and returns a
// StopResponse.
func Status(request models.StatusRequest,
	host string, permission Permission) (*models.StatusResponse, error) {

	httpRequest, err := StatusRequest(request, host)
	if err != nil {
		return nil,
			fmt.Errorf("error forming status request:\n %s", err.Error())
	}

	responseBody, err := send(httpRequest, permission)
	if err != nil {
		return nil,
			fmt.Errorf("error sending status request:\n %s", err.Error())
	}

	var response *models.StatusResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil,
			fmt.Errorf("error on unmarshal status response:\n %s", err.Error())
	}

	return response, nil
}

// generic send function to initialize a client and return the byte body of the
// response.
func send(request *http.Request, permission Permission) ([]byte, error) {
	client, err := newTLSClient(permission)
	if err != nil {
		return nil, fmt.Errorf("error initializing client:\n %s", err.Error())
	}

	responseBody, statusCode, err := client.SendRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request:\n %s", err.Error())
	}

	if statusCode >= 400 {

		var errorMsg *models.ErrorMessage
		err = json.Unmarshal(responseBody, &errorMsg)
		if err != nil {
			return nil,
				fmt.Errorf("error on unmarshal error response:\n %s", err.Error())
		}

		return nil,
			fmt.Errorf(
				"request returned with error\n status code: %d\n error message: %s",
				statusCode, errorMsg.Error)

	}

	return responseBody, nil
}
