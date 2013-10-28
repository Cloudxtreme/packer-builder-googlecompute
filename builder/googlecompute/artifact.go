package googlecompute

import (
	"fmt"
	"log"
)

type Artifact struct {
	imageName string
	client    *GoogleComputeClient
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying image: %s", a.imageName)
	operation, err := a.client.DeleteImage(a.imageName)
	if err != nil {
		return err
	}
	log.Print("Waiting for the instance to be deleted")
	for {
		status, err := a.client.GlobalOperationStatus(operation.Name)
		if err != nil {
			return err
		}
		if status == "DONE" {
			break
		}
	}
	return nil
}

func (*Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.imageName
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A disk image was created: %v", a.imageName)
}
