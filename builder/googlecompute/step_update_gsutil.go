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
type stepUpdateGsutil int

// Run executes the Packer build step that creates a GCE machine image.
func (s *stepUpdateGsutil) Run(state multistep.StateBag) multistep.StepAction {
	var (
		comm   = state.Get("communicator").(packer.Communicator)
		ui     = state.Get("ui").(packer.Ui)
	)
	ui.Say("Updating gsutil...")
	// Update gsutil now to prevent image creation from hanging.
	// The gcimagebundle command used in the create_image step calls gsutil
	// in the background, which can hang the image creation process with a prompt
	// to update gsutil if not running the latest version.
	cmd := new(packer.RemoteCmd)
	cmd.Command = "/usr/local/bin/gsutil update -n -f"
	err := cmd.StartWithUi(comm, ui)
	if err != nil {
		err := fmt.Errorf("Error updating gsutil: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (s *stepUpdateGsutil) Cleanup(state multistep.StateBag) {}
