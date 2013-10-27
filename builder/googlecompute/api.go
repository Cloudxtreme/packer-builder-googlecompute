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

// getZone returns a *compute.Zone representing the named zone.
// It returns an error if any.
func (g *GoogleComputeClient) getZone(name string) (*compute.Zone, error) {
	zoneGetCall := g.Service.Zones.Get(g.ProjectID, name)
	zone, err := zoneGetCall.Do()
	if err != nil {
		return nil, err
	}
	return zone, nil
}

// getMachineType returns a *compute.MachineType representing the named machine type.
// It returns an error if any.
func (g *GoogleComputeClient) getMachineType(name, zone string) (*compute.MachineType, error) {
	machineTypesGetCall := g.Service.MachineTypes.Get(g.ProjectID, zone, name)
	machineType, err := machineTypesGetCall.Do()
	if err != nil {
		return nil, err
	}
	if machineType.Deprecated == nil {
		return machineType, nil
	}
	return nil, errors.New("Machine Type does not exist: " + name)
}

// getImage returns a *compute.Image representing the named image.
// It returns an error if any.
func (g *GoogleComputeClient) getImage(name string) (*compute.Image, error) {
	// First try and find the image in the users project
	imagesGetCall := g.Service.ImagesService.Get(g.ProjectId, name)
	image, err := imagesGetCall.Do()
	if err != nil {
		log.Printf("Cannot find image: %s in project %s", name, g.ProjectId)
	}
	// Now try and find the image in the debian-cloud
	imagesGetCall := g.Service.ImagesService.Get("debian-cloud", name)
	image, err := imagesGetCall.Do()
	if err != nil {
		log.Printf("Cannot find image: %s in project %s", name, g.ProjectId)
	}
	if image != nil {
		if image.SelfLink != "" {
			return image, nil
		}
	}
	return nil, errors.New("Image does not exist: " + name)
}

// getNetwork returns a *compute.Network representing the named network.
// It returns an error if any.
func (g *GoogleComputeClient) getNetwork(name string) (*compute.Network, error) {
	networkGetCall := g.Service.NetworkService.Get(g.ProjectId, name)
	network, err := networkGetCall.Do()
	if err != nil {
		return nil, err
	}
	return network, nil
}

// InstanceConfig
type InstanceConfig struct {
	Description       string
	Image             string
	MachineType       string
	Metadata          *compute.Metadata
	Name              string
	NetworkInterfaces []*compute.NetworkInterface
	ServiceAccounts   []*compute.ServiceAccount
	Tags              []*compute.Tags
}

// createInstance.
func (g *GoogleComputeClient) createInstance(zone string, instanceConfig *InstanceConfig) (*compute.Operation, error) {
	// Attache disk
	instance := &compute.Instance{
		Description:       instanceConfig.Description,
		Image:             instanceConfig.Image,
		MachineType:       instanceConfig.MachineType,
		Metadata:          instanceConfig.Metadata,
		Name:              instanceConfig.Name,
		NetworkInterfaces: instanceConfig.NetworkInterfaces,
		ServiceAccounts:   instanceConfig.ServiceAccounts,
		Tags:              instanceConfig.Tags,
	}
	instanceInsertCall := g.Service.InstanceService.Insert(g.ProjectId, zone, instance)
	operation, err := instanceInsertCall.Do()
	if err != nil {
		return nil, err
	}
	return operation, nil
}

// sliceToTags converts a slice of strings to a *compute.Tags.
func sliceToTags(tags []string) *compute.Tags {
	return &compute.Tags{
		Items: tags,
	}
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
