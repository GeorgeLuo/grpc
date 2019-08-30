package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/GeorgeLuo/grpc/models"
)

// TODO persistent connection

var client *http.Client

func main() {

	if len(os.Args) < 2 {
		fmt.Println("start, stop, or command instruction not provided")
		os.Exit(1)
	}

	args := os.Args
	var certFile string
	var keyFile string
	var host string
	var endpoint string

	body, err := ReadArgs(args, &endpoint, &certFile, &keyFile, &host)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var request *http.Request

	switch endpoint {
	case "start":
		request, err = StartRequest(body, host)
	case "stop":
		request, err = StopRequest(body, host)
	case "status":
		request, err = StatusRequest(body, host)
	default:
		err = errors.New("invalid command")
		os.Exit(1)
	}

	client, err = newClient(certFile, keyFile, certFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	r, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer r.Body.Close()
	responseBody, responseError := ioutil.ReadAll(r.Body)
	if responseError != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s", responseBody)
}

// newClient creates a new tls client from key and certs
func newClient(certFile string, keyFile string, caCertFile string) (*http.Client, error) {

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// TODO add command line arg for cacert
	caCert, err := ioutil.ReadFile(caCertFile)
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

// TODO could return a config object holding this information

// ReadArgs populates the common parameters between the 3 endpoints and returns request interface.
func ReadArgs(args []string, endpoint *string, certFile *string, keyFile *string, host *string) (body interface{}, err error) {

	commonArgs := flag.NewFlagSet("api", flag.ExitOnError)

	certFilePtr := commonArgs.String("cert", "", "path to cert file")
	keyFilePtr := commonArgs.String("key", "", "path to key file")
	hostPtr := commonArgs.String("host", "", "endpoint of request")
	*endpoint = args[1]

	switch *endpoint {
	case "start":
		startCommandPtr := commonArgs.String("command", "", "command to exec")
		commonArgs.Parse(os.Args[2:])
		body = models.StartRequest{Command: *startCommandPtr}
	case "stop":
		stopTaskID := commonArgs.String("task_id", "", "task_id of process")
		commonArgs.Parse(os.Args[2:])
		body = models.StopRequest{TaskID: *stopTaskID}
	case "status":
		statusTaskID := commonArgs.String("task_id", "", "task_id of process")
		commonArgs.Parse(os.Args[2:])
		body = models.StatusRequest{TaskID: *statusTaskID}
	default:
		err = errors.New("invalid command")
	}

	*certFile = *certFilePtr
	*keyFile = *keyFilePtr
	*host = *hostPtr

	return body, err
}
