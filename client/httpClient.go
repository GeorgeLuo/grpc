package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

// Permission encapsulates the tls configurations.
type Permission struct {
	CertFile   string
	KeyFile    string
	CaCertFile string
}

// GetCaCert returns the CertFile if CaCert is empty.
func (perm *Permission) GetCaCert() string {
	if perm.CaCertFile == "" {
		return perm.CertFile
	}
	return perm.CaCertFile
}

// TLSClient is an extension of the go default http client. This abstraction
// exists for future changes when the underlying client may need to be replaced
// for additional features
type TLSClient struct {
	client http.Client
}

// newTLSClient creates a new tls client from key and certs provided by
// permission object
func newTLSClient(permission Permission) (*TLSClient, error) {

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

	return &TLSClient{
		client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{cert},
				},
			},
		},
	}, nil
}

// SendRequest encapsulates the entire process of initializing the client
// and sending a request, returning the byte body of the response.
func (client *TLSClient) SendRequest(request *http.Request) ([]byte, error) {

	var err error

	r, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	responseBody, responseError := ioutil.ReadAll(r.Body)
	if responseError != nil {
		return nil, responseError
	}

	return responseBody, nil
}
