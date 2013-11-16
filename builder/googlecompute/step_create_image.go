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
	imageBundleCmd := "/usr/share/imagebundle/image_bundle.py -r / -o /tmp/"
	cmd := new(packer.RemoteCmd)
	cmd.Command = fmt.Sprintf("%s --output_file_name %s -b %s", imageBundleCmd, c.ImageName, c.BucketName)
	ui.Say("Creating image using: " + cmd.Command)
	err := cmd.StartWithUi(comm, ui)
	if err != nil {
		ui.Error(fmt.Sprintf("Error creating image"))
	}
	state.Put("image_name", "random-image")
	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
	return
}
