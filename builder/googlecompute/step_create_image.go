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
	ui.Say("Creating image using /usr/share/imagebundle/image_bundle.py")
	//c := state.Get("config").(config)
	cmd := new(packer.RemoteCmd)
	cmd.Command = "/usr/share/imagebundle/image_bundle.py -r / -o /tmp/"
	err := cmd.StartWithUi(comm, ui)
	if err != nil {
		ui.Error(fmt.Sprintf("Error creating image"))
	}
	ui.Say("Uploading image using gsutil")
	cmd.Command = "/usr/local/bin/gsutil cp /tmp/*.image.tar.gz gs://packer-images"
	err = cmd.StartWithUi(comm, ui)
	if err != nil {
		ui.Error(fmt.Sprintf("Error creating image"))
	}
	state.Put("image_name", "random-image")
	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
	return
}
