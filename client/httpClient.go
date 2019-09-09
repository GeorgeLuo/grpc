package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var client *http.Client

// Permission encapsulates the tls configurations.
type Permission struct {
	CertFile string
	KeyFile  string
	CaCert   string
}

// GetCaCert returns the CertFile if CaCert is empty.
func (perm *Permission) GetCaCert() string {
	if perm.CaCert == "" {
		return perm.CertFile
	}
	return perm.CaCert
}

// SendRequest encapsulates the entire process of initializing the client
// and sending a request, returning the byte body of the response.
func SendRequest(permission Permission,
	// request *http.Request, target interface{}) (*[]byte, error) {
	request *http.Request, target interface{}) error {

	var err error

	if client == nil {
		client, err = newClient(permission)
		if err != nil {
			// return nil, err
			return err
		}
	}

	r, err := client.Do(request)
	if err != nil {
		// return nil, err
		return err
	}

	defer r.Body.Close()
	responseBody, responseError := ioutil.ReadAll(r.Body)
	if responseError != nil {
		// return nil, responseError
		return responseError
	}
	return json.Unmarshal(responseBody, target)
	// return &responseBody, nil
}

// newClient creates a new tls client from key and certs
func newClient(permission Permission) (*http.Client, error) {

	cert, err := tls.LoadX509KeyPair(permission.CertFile, permission.KeyFile)
	if err != nil {
		return nil, err
	}

	// TODO add command line arg for cacert
	caCert, err := ioutil.ReadFile(permission.GetCaCert())
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}, nil
}