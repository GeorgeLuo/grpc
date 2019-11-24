package utils

import (
	"encoding/json"
	"os"
)

// JSON config model

// LoadServerConfiguration takes a string filename and returns a Server Config
// object populated with file's values
func LoadServerConfiguration(file string) (ServerConfig, error) {
	var config ServerConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return DefaultServerConfig(), err
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, nil
}

// ServerConfig defines the hosting configurations of a remote server.
type ServerConfig struct {
	// a Network is a literal endpoint of of an execution network to join.
	NetworkConf struct {
		// Host is the endpoint address of the network.
		Host string `json:"host,omitempty"`
		// Alias is the is the key to a mapping to a host network.
		Alias string `json:"alias,omitempty"`
	} `json:"network"`
	Port string `json:"port"`
	// WorkingDir is the relative start directory (relevant for bin).
	WorkingDir string `json:"workingDir"`
	// LogDir is where process output will be written.
	LogDir     string `json:"logDir"`
	CertFile   string `json:"certFile"`
	CaCertFile string `json:"cacertFile"`
	KeyFile    string `json:"keyFile"`
}

// DefaultServerConfig is used to return an empty SyncMap.
func DefaultServerConfig() ServerConfig {
	serverConfig := ServerConfig{
		Port: "8443",
		// default to locations where the server is called from.
		WorkingDir: "./",
		LogDir:     "./",
		CertFile:   "cert.pem",
		CaCertFile: "cert.pem",
		KeyFile:    "key.pem",
	}
	return serverConfig
}
