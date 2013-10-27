package googlecompute

import (
	"errors"
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

// New return a new GoogleComputeClient.
func New(c *ClientSecrets, pemKey []byte) (*GoogleComputeClient, error) {
	googleComputeClient := &GoogleComputeClient{}
	googleComputeClient.ProjectId = extractProjectId(c.Web.ClientId)
	// Get the access token.
	t := jwt.NewToken(c.Web.ClientEmail, scopes(), pemKey)
	t.ClaimSet.Aud = c.Web.TokenURI
	httpClient := &http.Client{}
	token, err := t.Assert(httpClient)
	if err != nil {
		return nil, err
	}
	// Create the Google Compute client.
	config := &oauth.Config{
		ClientId: c.Web.ClientId,
		Scope:    scopes(),
		TokenURL: c.Web.TokenURI,
		AuthURL:  c.Web.AuthURI,
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
	zoneGetCall := g.Service.Zones.Get(g.ProjectId, name)
	zone, err := zoneGetCall.Do()
	if err != nil {
		return nil, err
	}
	return zone, nil
}

// getMachineType returns a *compute.MachineType representing the named machine type.
// It returns an error if any.
func (g *GoogleComputeClient) getMachineType(name, zone string) (*compute.MachineType, error) {
	machineTypesGetCall := g.Service.MachineTypes.Get(g.ProjectId, zone, name)
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
	imagesGetCall := g.Service.Images.Get(g.ProjectId, name)
	image, err := imagesGetCall.Do()
	if err != nil {
		log.Printf("Cannot find image: %s in project %s", name, g.ProjectId)
	}
	// Now try and find the image in the debian-cloud
	imagesGetCall = g.Service.Images.Get("debian-cloud", name)
	image, err = imagesGetCall.Do()
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
	networkGetCall := g.Service.Networks.Get(g.ProjectId, name)
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
	Tags              *compute.Tags
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
		Tags:              instanceConfig.Tags,
	}
	instanceInsertCall := g.Service.Instances.Insert(g.ProjectId, zone, instance)
	operation, err := instanceInsertCall.Do()
	if err != nil {
		return nil, err
	}
	return operation, nil
}

// NewNetworkInterface returns a *compute.NetworkInterface.
func NewNetworkInterface(network *compute.Network, public bool) *compute.NetworkInterface {
	accessConfigs := make([]*compute.AccessConfig, 0)
	if public {
		c := &compute.AccessConfig{
			Name: "AccessConfig created by Packer",
			Type: "ONE_TO_ONE_NAT",
		}
		accessConfigs = append(accessConfigs, c)
	}
	return &compute.NetworkInterface{
		AccessConfigs: accessConfigs,
		Network:       network.SelfLink,
	}
}

// mapToMetadata converts a map[string]string to a *compute.Metadata.
func mapToMetadata(metadata map[string]string) *compute.Metadata {
	items := make([]*compute.MetadataItems, len(metadata))
	for k, v := range metadata {
		items = append(items, &compute.MetadataItems{k, v})
	}
	return &compute.Metadata{
		Items: items,
	}
}

// sliceToTags converts a []string to a *compute.Tags.
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
