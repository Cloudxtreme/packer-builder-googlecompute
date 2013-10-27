package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/goauth2/oauth/jwt"
	"code.google.com/p/google-api-go-client/compute/v1beta16"
)

type MachineMap map[string]*MachineType

type MachineType struct {
	Name  string
	State string
	Url   string
	Zone  string
}

type ZoneMap map[string]*Zone

type Zone struct {
	Name string
	Url  string
}

type ClientSecrets struct {
	Web struct {
		ClientEmail string `json:"client_email"`
		ClientID    string `json:"client_id"`
		AuthURI     string `json:"auth_uri"`
		TokenURI    string `json:"token_uri"`
	}
}

var clientSecrets ClientSecrets
var zone = "us-central1-a"

func main() {
	scopeList := []string{
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/compute.readonly",
	}
	scope := strings.Join(scopeList, " ")
	// read in the key
	pemKeyBytes, err := ioutil.ReadFile("gce-service.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Craft the ClaimSet and JWT token.
	secretBytes, err := ioutil.ReadFile("client_secrets.json")
	if err != nil {
		log.Fatalf("Cannot load client secrets file: %s", err.Error())
	}
	err = json.Unmarshal(secretBytes, &clientSecrets)
	if err != nil {
		log.Fatalf("Cannot unmarshal client secrets: %s", err.Error())
	}
	projectID := strings.SplitN(clientSecrets.Web.ClientID, "-", 2)[0]
	t := jwt.NewToken(clientSecrets.Web.ClientEmail, scope, pemKeyBytes)
	t.ClaimSet.Aud = clientSecrets.Web.TokenURI
	c := &http.Client{}

	// Get the access token.
	token, err := t.Assert(c)
	if err != nil {
		log.Fatal("Failed to get token", err.Error())
	}

	// Make the request
	config := &oauth.Config{
		ClientId:   clientSecrets.Web.ClientID,
		Scope:      scope,
		TokenURL:   clientSecrets.Web.TokenURI,
		AuthURL:    clientSecrets.Web.AuthURI,
		TokenCache: oauth.CacheFile("cache.json"),
	}
	transport := &oauth.Transport{Config: config}
	transport.Token = token

	s, err := compute.New(transport.Client())
	if err != nil {
		log.Fatal("Cannot create compute service", err.Error())
	}

	// Create the zoneMap
	zoneMap := make(ZoneMap)
	zoneListCall := s.Zones.List(projectID)
	zoneList, err := zoneListCall.Do()
	if err != nil {
		log.Fatal("Failed on zonelist call", err.Error())
	}
	for _, z := range zoneList.Items {
		zoneMap[z.Name] = &Zone{
			Name: z.Name,
			Url:  z.SelfLink,
		}
	}

	if _, ok := zoneMap[zone]; !ok {
		log.Fatalf("Zone: %s does not exist", zone)
	}

	// Create the Machine Map
	machineMap := make(MachineMap)
	machineTypesCall := s.MachineTypes.List(projectID, zone)
	machineTypeList, err := machineTypesCall.Do()
	if err != nil {
		log.Fatal("Error gathering machine types: ", err.Error())
	}
	for _, mt := range machineTypeList.Items {
		machineMap[mt.Name] = &MachineType{
			Name: mt.Name,
			Url:  mt.SelfLink,
			Zone: mt.Zone,
		}
		if mt.Deprecated != nil {
			machineMap[mt.Name].State = mt.Deprecated.State
		}
	}

	// List all the images
	imagesListCall := s.Images.List(projectID)
	imageList, err := imagesListCall.Do()
	if err != nil {
		log.Fatal("Failed on images list call", err.Error())
	}
	for _, i := range imageList.Items {
		fmt.Println(i.Name)
		fmt.Println(i.SelfLink)
	}
}
