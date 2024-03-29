package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/GeorgeLuo/grpc/models"
)

// TODO persistent connection
// TODO: add client test code for setting up tls

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
	startAlias := startCommand.String("alias", "", "alias for process")
	startFormatted := startCommand.Bool("t", false, "generate output as table")

	stopCommand := flag.NewFlagSet("stop", flag.ExitOnError)
	stopCertFile := stopCommand.String("cert", "cert.pem", "path to cert file")
	stopCaCertFile := stopCommand.String("cacert", "", "path to cacert file")
	stopKeyFile := stopCommand.String("key", "key.pem", "path to key file")
	stopHost := stopCommand.String("host", "localhost", "endpoint of request")
	stopTaskID := stopCommand.String("task_id", "", "task_id of process")
	stopAlias := stopCommand.String("alias", "", "alias for process")
	stopFormatted := stopCommand.Bool("t", false, "generate output as table")

	var batchTaskID arrayFlags

	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)
	statusCertFile := statusCommand.String("cert", "cert.pem", "path to cert file")
	statusCaCertFile := statusCommand.String("cacert", "", "path to cacert file")
	statusKeyFile := statusCommand.String("key", "key.pem", "path to key file")
	statusHost := statusCommand.String("host", "localhost", "endpoint of request")
	statusAlias := statusCommand.String("alias", "", "alias for process")
	statusFormatted := statusCommand.Bool("t", false, "generate output as table")

	statusCommand.Var(&batchTaskID, "task_id", "task_id of process")

	var err error

	switch os.Args[1] {
	case "start":
		err = startCommand.Parse(os.Args[2:])
	case "stop":
		err = stopCommand.Parse(os.Args[2:])
	case "status":
		err = statusCommand.Parse(os.Args[2:])
	default:
		fmt.Println("invalid command")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("error parsing request:\n %s\n", err.Error())
		os.Exit(1)
	}

	var permission Permission
	var renderable []Renderable
	var tabled bool

	if startCommand.Parsed() {
		permission = Permission{
			CertFile:   *startCertFile,
			KeyFile:    *startKeyFile,
			CaCertFile: *startCaCertFile,
		}

		startResponse, err := Start(models.StartRequest{
			Command: *startExec,
			Alias:   *startAlias,
		}, *startHost, permission)

		if err != nil {
			fmt.Printf("error sending start:\n %s\n", err.Error())
			os.Exit(1)
		}

		renderable = append(renderable, NewRenderableStartResponse(startResponse))
		tabled = *startFormatted

	} else if stopCommand.Parsed() {
		permission = Permission{
			CertFile:   *stopCertFile,
			KeyFile:    *stopKeyFile,
			CaCertFile: *stopCaCertFile,
		}

		if (*stopAlias != "") == (*stopTaskID != "") {
			fmt.Println("error: must provide one (and only one) of task_id or alias")
			os.Exit(1)
		}

		var stopResponse *models.StopResponse
		if *stopAlias != "" {
			stopResponse, err = Stop(models.StopRequest{Alias: *stopAlias},
				*stopHost, permission)
		} else {
			stopResponse, err = Stop(models.StopRequest{TaskID: *stopTaskID},
				*stopHost, permission)
		}

		if err != nil {
			fmt.Printf("error sending stop:\n %s\n", err.Error())
			os.Exit(1)
		}

		renderable = append(renderable, NewRenderableStopResponse(stopResponse))
		tabled = *stopFormatted

	} else if statusCommand.Parsed() {

		permission = Permission{
			CertFile:   *statusCertFile,
			KeyFile:    *statusKeyFile,
			CaCertFile: *statusCaCertFile,
		}

		if *statusAlias != "" && len(batchTaskID) > 0 {
			fmt.Println("error: provide either alias or task_id, not both")
			os.Exit(1)
		}

		// TODO: refactor this from main()

		// only evaluate alias if alias is present. TODO: The renderable being added
		// will eventually contain multiple processes. The Status method will return
		// a slice of renderable responses.
		if *statusAlias != "" {
			statusResponse, err := Status(models.StatusRequest{
				Alias: *statusAlias,
			}, *statusHost, permission)

			if err != nil {
				fmt.Printf("error getting status for alias %s:\n %s\n",
					*statusAlias, err.Error())
			} else {
				b := NewBatchRenderable(*statusAlias)
				b.AddRow(NewRenderableStatusResponse(statusResponse))
				renderable = append(renderable, b)
			}
		}

		b := NewBatchRenderable("")

		// TODO: return error on task_id not provided
		for _, statusTaskID := range batchTaskID {

			statusResponse, err := Status(models.StatusRequest{
				TaskID: statusTaskID,
			}, *statusHost, permission)

			if err != nil {
				fmt.Printf("error getting status for task %s:\n %s\n",
					statusTaskID, err.Error())
			} else {
				b.AddRow(NewRenderableStatusResponse(statusResponse))
			}
		}

		if b.Size() > 0 {
			renderable = append(renderable, b)
		}
		tabled = *statusFormatted

	} else {
		fmt.Println("error parsing arguments")
		os.Exit(1)
	}

	Render(os.Stdout, renderable, tabled)
}

// arrayFlags is used to manage batch requests to status using multiple task_ids
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
