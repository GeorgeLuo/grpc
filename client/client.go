package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/GeorgeLuo/grpc/models"
)

// TODO persistent connection
// TODO: add client test code for setting up tls

type arrayFlags []string

func (flags *arrayFlags) String() string {
	var sb strings.Builder
	first := true
	for _, flag := range *flags {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(flag)
		first = false
	}

	return sb.String()
}

func (flags *arrayFlags) Set(value string) error {
	*flags = append(*flags, value)
	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("start, stop, or command instruction not provided")
		os.Exit(1)
	}

	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	startCertFile := startCommand.String("cert", "cert.pem", "path to cert file")
	startCaCertFile := startCommand.String("cacert", "", "path to cacert file")
	startKeyFile := startCommand.String("key", "key.pem", "path to key file")
	startHost := startCommand.String("host", "localhost", "endpoint of request")
	startExec := startCommand.String("command", "", "command to exec")

	stopCommand := flag.NewFlagSet("stop", flag.ExitOnError)
	stopCertFile := stopCommand.String("cert", "cert.pem", "path to cert file")
	stopCaCertFile := stopCommand.String("cacert", "", "path to cacert file")
	stopKeyFile := stopCommand.String("key", "key.pem", "path to key file")
	stopHost := stopCommand.String("host", "localhost", "endpoint of request")
	stopTaskID := stopCommand.String("task_id", "", "task_id of process")

	var batchTaskID arrayFlags

	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)
	statusCertFile := statusCommand.String("cert", "cert.pem", "path to cert file")
	statusCaCertFile := statusCommand.String("cacert", "", "path to cacert file")
	statusKeyFile := statusCommand.String("key", "key.pem", "path to key file")
	statusHost := statusCommand.String("host", "localhost", "endpoint of request")
	statusCommand.Var(&batchTaskID, "task_id", "task_id of process")

	switch os.Args[1] {
	case "start":
		startCommand.Parse(os.Args[2:])
	case "stop":
		stopCommand.Parse(os.Args[2:])
	case "status":
		statusCommand.Parse(os.Args[2:])
	default:
		fmt.Println("invalid command")
		os.Exit(1)
	}

	var request *http.Request
	var permission Permission

	var err error
	var batchKey *string

	if startCommand.Parsed() {
		permission = Permission{*startCertFile, *startKeyFile, *startCaCertFile}
		request, err = StartRequest(models.StartRequest{Command: *startExec}, *startHost)
	} else if stopCommand.Parsed() {
		permission = Permission{*stopCertFile, *stopKeyFile, *stopCaCertFile}
		request, err = StopRequest(models.StopRequest{TaskID: *stopTaskID}, *stopHost)
	} else if statusCommand.Parsed() {
		permission = Permission{*statusCertFile, *statusKeyFile, *statusCaCertFile}
		if len(batchTaskID) > 1 {
			batchKey = new(string)
			*batchKey = "status_responses"
			request, err = StatusBatchRequest(models.StatusBatchRequest{TaskIDs: batchTaskID}, *statusHost)
		} else {
			request, err = StatusRequest(models.StatusRequest{TaskID: batchTaskID[0]}, *statusHost)
		}
	} else {
		fmt.Println("error parsing arguments")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("error forming request: [%s]\n", err.Error())
		os.Exit(1)
	}

	var responseBody *models.StatusResponse
	// var responseBody interface{}

	// responseBody, err := SendRequest(permission, request)
	err = SendRequest(permission, request, &responseBody)
	if err != nil {
		fmt.Printf("error sending request: [%s]\n", err.Error())
		os.Exit(1)
	}

	// PrintTable(batchKey, responseBody, os.Stdout)
	fmt.Println(responseBody)
}
