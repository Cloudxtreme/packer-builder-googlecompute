package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kelseyhightower/packer-builder-googlecompute/builder/googlecompute"
)

var (
	clientSecretsFile string
	privateKeyFile    string
)

func init() {
	flag.StringVar(&clientSecretsFile, "s", "", "path to client secrets file.")
	flag.StringVar(&privateKeyFile, "k", "", "path to private key file.")
}

func main() {
	flag.Parse()
	if clientSecretsFile == "" || privateKeyFile == "" {
		log.Fatal("-s and -k are required")
	}
	pemKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	clientSecretsBytes, err := ioutil.ReadFile(clientSecretsFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	var clientSecrets *googlecompute.ClientSecrets
	err = json.Unmarshal(clientSecretsBytes, &clientSecrets)
	if err != nil {
		log.Fatal(err.Error())
	}
}
