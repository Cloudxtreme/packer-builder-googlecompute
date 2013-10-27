package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"code.google.com/p/google-api-go-client/compute/v1beta16"
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
	var (
		imageName       = "debian-7-wheezy-v20130926"
		machineTypeName = "n1-standard-1-d"
		networkName     = "default"
		zoneName        = "us-central2-a"
		projectId       = "hightower-labs"
	)
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
	g, err := googlecompute.New(projectId, clientSecrets, pemKeyBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	instanceConfig := &googlecompute.InstanceConfig{
		Description: "New instance created by Packer",
		Name:        "packer-instance",
	}
	// Validate the zone.
	zone, err := g.GetZone(zoneName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get the image
	image, err := g.GetImage(imageName)
	if err != nil {
		log.Fatal(err.Error())
	}
	instanceConfig.Image = image.SelfLink

	machineType, err := g.GetMachineType(machineTypeName, zone.Name)
	if err != nil {
		log.Fatal(err.Error())
	}
	instanceConfig.MachineType = machineType.SelfLink

	network, err := g.GetNetwork(networkName)
	if err != nil {
		log.Fatal(err.Error())
	}
	networkInterface := googlecompute.NewNetworkInterface(network, true)
	networkInterfaces := []*compute.NetworkInterface{
		networkInterface,
	}
	instanceConfig.NetworkInterfaces = networkInterfaces

	operation, err := g.CreateInstance(zone.Name, instanceConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", operation)
}
