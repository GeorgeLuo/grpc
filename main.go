package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/GeorgeLuo/grpc/core"
	"github.com/GeorgeLuo/grpc/jobs"
	"github.com/GeorgeLuo/grpc/utils"
	"github.com/gorilla/mux"
)

// TODO: add a handler to eventstream output of process

// StaticLogger is the global static logger
var StaticLogger utils.StaticLogger

func main() {
	var config utils.ServerConfig
	var err error

	if len(os.Args) < 2 {
		log.Println("operation not provided")
		os.Exit(1)
	}

	startServer := flag.NewFlagSet("start", flag.ExitOnError)
	configFile := startServer.String("conf", "grpc_server_conf.json",
		"start configurations")

	switch os.Args[1] {
	case "start":
		err = startServer.Parse(os.Args[2:])
	default:
		fmt.Println("invalid command")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("error parsing request:\n %s\n", err.Error())
		os.Exit(1)
	}

	log.Printf("loading configurations from: %s", *configFile)
	config, err = utils.LoadServerConfiguration(*configFile)
	if err != nil {
		log.Printf("config file not not found, using default configurations")
	}

	prettyConf, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("error printing configs: ", err)
		os.Exit(1)
	}

	StaticLogger = utils.NewStaticLogger(config.LogDir, "grpc_server.log")
	StaticLogger.SetLevel(utils.D)
	StaticLogger.SetGlobalPrepend("grpc_server")
	core.StaticLogger = &StaticLogger
	defer StaticLogger.Close()

	StaticLogger.WriteDtInfo(string(prettyConf))
	StaticLogger.WriteDtInfo("setting up handlers ...")

	router := mux.NewRouter()
	router.HandleFunc("/status", core.StatusHandler).
		Methods("POST")
	router.HandleFunc("/start", core.StartHandler).
		Methods("POST")
	router.HandleFunc("/stop", core.StopHandler).
		Methods("POST")
	router.HandleFunc("/jobs/status", jobs.JobStatusHandler).
		Methods("POST")
	router.HandleFunc("/jobs/start", jobs.JobStartHandler).
		Methods("POST")
	router.HandleFunc("/jobs/status", jobs.JobStatusHandler).
		Methods("POST")
	router.HandleFunc("/jobs/stop", jobs.JobStopHandler).
		Methods("POST")

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile(config.CaCertFile)
	if err != nil {
		StaticLogger.WriteDtFatal(err.Error())
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Handler:   router,
		Addr:      "0.0.0.0:" + config.Port,
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS(config.CertFile, config.KeyFile))
}
