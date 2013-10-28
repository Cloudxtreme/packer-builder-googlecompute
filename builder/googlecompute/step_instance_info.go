package googlecompute

import (
	"github.com/mitchellh/multistep"
)

type stepInstanceInfo struct {
	ip string
}

func (s *stepInstanceInfo) Run(state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *stepInstanceInfo) Cleanup(state multistep.StateBag) {
	return
}
