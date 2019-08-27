package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// TODO persistent connection

func main() {

	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	stopCommand := flag.NewFlagSet("stop", flag.ExitOnError)
	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)

	startCommandPtr := startCommand.String("command", "", "command to exec")
	startCertPtr := startCommand.String("cert", "", "path to cert file")
	startKeyPtr := startCommand.String("key", "", "path to key file")
	startHostPtr := startCommand.String("host", "", "endpoint of request")

	// TODO refactor repetitive code
	stopTaskIDPtr := stopCommand.String("task_id", "", "command to exec")
	stopCertPtr := stopCommand.String("cert", "", "path to cert file")
	stopKeyPtr := stopCommand.String("key", "", "path to key file")
	stopHostPtr := stopCommand.String("host", "", "endpoint of request")

	statusTaskIDPtr := statusCommand.String("task_id", "", "command to exec")
	statusCertPtr := statusCommand.String("cert", "", "path to cert file")
	statusKeyPtr := statusCommand.String("key", "", "path to key file")
	statusHostPtr := statusCommand.String("host", "", "endpoint of request")

	if len(os.Args) < 2 {
		fmt.Println("start, stop, or command instruction not provided")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCommand.Parse(os.Args[2:])
	case "stop":
		stopCommand.Parse(os.Args[2:])
	case "status":
		statusCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	var certFile string
	var keyFile string

	var request *http.Request
	var err error

	if startCommand.Parsed() {
		if *startCommandPtr == "" {
			startCommand.PrintDefaults()
			os.Exit(1)
		}

		certFile = *startCertPtr
		keyFile = *startKeyPtr

		var data = []byte(`{"command":"` + *startCommandPtr + `"}`)

		urlString := "https://" + *startHostPtr + ":8443/start"

		request, err = http.NewRequest("POST", urlString, bytes.NewBuffer(data))
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			os.Exit(1)
		}
	}

	if stopCommand.Parsed() {
		if *stopTaskIDPtr == "" {
			stopCommand.PrintDefaults()
			os.Exit(1)
		}

		certFile = *stopCertPtr
		keyFile = *stopKeyPtr

		var data = []byte(`{"task_id":"` + *stopTaskIDPtr + `"}`)
		urlString := "https://" + *stopHostPtr + ":8443/stop"

		request, err = http.NewRequest("POST", urlString, bytes.NewBuffer(data))
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			os.Exit(1)
		}
	}

	if statusCommand.Parsed() {
		if *statusTaskIDPtr == "" {
			statusCommand.PrintDefaults()
			os.Exit(1)
		}

		certFile = *statusCertPtr
		keyFile = *statusKeyPtr
		urlString := "https://" + *statusHostPtr + ":8443/status/" + *statusTaskIDPtr

		request, err = http.NewRequest("GET", urlString, nil)
		if err != nil {
			os.Exit(1)
		}
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	// Request /hello via the created HTTPS client over port 8443 via GET
	r, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body to stdout
	fmt.Printf("%s", body)
}

// TODO refactor handlers for instructions to return request objects
