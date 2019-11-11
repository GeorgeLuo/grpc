package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GeorgeLuo/grpc/core"
	"github.com/GeorgeLuo/grpc/jobs"
	"github.com/gorilla/mux"
)

// TODO: add a handler to eventstream output of process

func main() {
	log.Println("setting up handlers ...")
	router := mux.NewRouter()
	router.HandleFunc("/status", core.StatusHandler).
		Methods("POST")
	router.HandleFunc("/start", core.StartHandler).
		Methods("POST")
	router.HandleFunc("/stop", core.StopHandler).
		Methods("POST")
	router.HandleFunc("/job/status", jobs.JobStatusHandler).
		Methods("POST")
	router.HandleFunc("/job/start", jobs.JobStartHandler).
		Methods("POST")

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
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
		Addr:      "0.0.0.0:8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
