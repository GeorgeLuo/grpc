package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

func main() {
	fmt.Println("setting up handlers ...")
	router := mux.NewRouter()
	router.HandleFunc("/status/{task_id}", StatusHandler).
		Methods("GET")
	router.HandleFunc("/start", StartHandler).
		Methods("POST")
	router.HandleFunc("/stop", StopHandler).
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
		ClientCAs: caCertPool,
		// ClientAuth: tls.NoClientCert,
		// InsecureSkipVerify: true,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Handler: router,
		Addr:      "0.0.0.0:8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait

	// log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
	log.Fatal(server.ListenAndServe())
}
