package googlecompute

import (
	"github.com/mitchellh/multistep"
)

type stepCreateImage struct {
	imageName string
}

func (s *stepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	state.Put("image_name", "random-image")
	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {
	return
}
