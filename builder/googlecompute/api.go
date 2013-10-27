package googlecompute

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

// ClientSecrets represents the parsed client secrets of a Google Compute
// service account.
type ClientSecrets struct {
	Web struct {
		ClientEmail string `json:"client_email"`
		ClientId    string `json:"client_id"`
		AuthURI     string `json:"auth_uri"`
		TokenURI    string `json:"token_uri"`
	}
}

// GoogleComputeClient represents a Google Compute Engine client.
type GoogleComputeClient struct {
	ProjectId     string
	Service       *compute.Service
	clientSecrets *ClientSecrets
}

// NewGoogleComputeClient return a new GoogleComputeClient.
func NewGoogleComputeClient(c *ClientSecrets, pemKey []byte) (*GoogleComputeClient, error) {
	googleComputeClient := &GoogleComputeClient{}
	googleComputeClient.ProjectId = extractProjectId(c.Web.ClientID)
	// Get the access token.
	t := jwt.NewToken(c.Web.ClientEmail, scopes(), pemKey)
	t.ClaimSet.Aud = c.Web.TokenURI
	c := &http.Client{}
	token, err := t.Assert(c)
	if err != nil {
		return nil, err
	}
	// Create the Google Compute client.
	config := &oauth.Config{
		ClientId: clientSecrets.Web.ClientID,
		Scope:    scope,
		TokenURL: clientSecrets.Web.TokenURI,
		AuthURL:  clientSecrets.Web.AuthURI,
	}
	transport := &oauth.Transport{Config: config}
	transport.Token = token

	s, err := compute.New(transport.Client())
	if err != nil {
		return nil, err
	}
	googleComputeClient.Service = s
	return googleComputeClient, nil
}

// getZoneURL returns the fully-qualified URL of the named zone resource.
// It returns an error if any.
func (g *GoogleComputeClient) getZoneURL(name string) (string, error) {
	zoneListCall := g.Service.Zones.List(g.ProjectID)
	zoneList, err := zoneListCall.Do()
	if err != nil {
		return "", err
	}
	for _, z := range zoneList.Items {
		if z.Name == name {
			return z.SelfLink, nil
		}
	}
	return "", errors.New("Zone does not exits: " + name)
}

// getMachineTypeURL returns the fully-qualified URL of the named machine type.
// It returns an error if any.
func (g *GoogleComputeClient) getMachineTypeURL(name, zone string) (string, error) {

	machineTypesCall := g.Service.MachineTypes.List(g.ProjectID, zone)
	machineTypeList, err := machineTypesCall.Do()
	if err != nil {
		return "", err
	}
	for _, mt := range machineTypeList.Items {
		if mt.Name == name && mt.Deprecated == nil {
			return mt.SelfLink, nil
		}
	}
	return "", errors.New("Machine Type does not exist: " + name)
}

// scopes return a space separated list of scopes.
func scopes() string {
	s := []string{
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/compute.readonly",
		"https://www.googleapis.com/auth/devstorage.full_control",
		"https://www.googleapis.com/auth/devstorage.read_write",
		"https://www.googleapis.com/auth/devstorage.write_only",
	}
	return strings.Join(s, " ")
}

// extractProjectId returns a string representing the Project ID.
func extractProjectId(clientId string) string {
	return strings.SplitN(clientId, "-", 2)[0]
}
