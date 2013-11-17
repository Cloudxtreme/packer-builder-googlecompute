package googlecompute

import (
	"fmt"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepCreateImage struct {
	imageName string
}

func (s *stepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	comm := state.Get("communicator").(packer.Communicator)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)
	client := state.Get("client").(*GoogleComputeClient)
	imageBundleCmd := "/usr/share/imagebundle/image_bundle.py -r / -o /tmp/"
	cmd := new(packer.RemoteCmd)
	cmd.Command = fmt.Sprintf("%s --output_file_name %s -b %s", imageBundleCmd, c.ImageName, c.BucketName)
	ui.Say("Creating image using: " + cmd.Command)
	err := cmd.StartWithUi(comm, ui)
	if err != nil {
		err := fmt.Errorf("Error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	imageURL := fmt.Sprintf("https://storage.cloud.google.com/%s/%s", c.BucketName, c.ImageName)
	operation, err := client.CreateImage(c.ImageName, c.ImageDescription, imageURL)
	if err != nil {
		err := fmt.Errorf("Error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("Waiting for image to become available...")
	err = waitForGlobalOperationState("DONE", operation.Name, client, c.stateTimeout)
	if err != nil {
		err := fmt.Errorf("Error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("image_name", c.ImageName)
	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
	return
}
