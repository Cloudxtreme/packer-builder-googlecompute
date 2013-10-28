package googlecompute

import (
	"github.com/mitchellh/multistep"
)

type stepCreateSSHKey struct {
	instanceName uint
}

func (s *stepCreateSSHKey) Run(state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	return
}
