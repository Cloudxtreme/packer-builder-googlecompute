// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package googlecompute

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// stepCreateImage represents a Packer build step that creates GCE machine images.
type stepCreateImage struct {
	imageName string
}

// Run executes the Packer build step that creates a GCE machine image.
func (s *stepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	var (
		config = state.Get("config").(config)
		comm   = state.Get("communicator").(packer.Communicator)
		ui     = state.Get("ui").(packer.Ui)
	)
	ui.Say("Creating image...")
	// Google Compute images must be created using the image_bundle.py utility
	// from the target GCE instance. Next the image must be uploaded to a Google
	// Cloud Storage bucket before it can be made available to the GCE project.
	imageBundleCmd := "/usr/share/imagebundle/image_bundle.py -r / -o /tmp/"
	cmd := new(packer.RemoteCmd)
	cmd.Command = fmt.Sprintf("%s --output_file_name %s.tar.gz -b %s",
		imageBundleCmd, config.ImageName, config.BucketName)
	err := cmd.StartWithUi(comm, ui)
	if err != nil {
		err := fmt.Errorf("Error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {}
