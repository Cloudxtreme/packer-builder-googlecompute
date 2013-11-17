// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

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
	// Ignore the operation result as we are not waiting until it completes.
	_, err := a.client.DeleteImage(a.imageName)
	if err != nil {
		return err
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
