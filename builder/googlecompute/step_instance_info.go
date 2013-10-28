package googlecompute

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepInstanceInfo struct{}

func (s *stepInstanceInfo) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*GoogleComputeClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)
	instanceName := state.Get("instance_name").(string)
	instanceOperationName := state.Get("instance_operation_name").(string)
	ui.Say("Waiting for the instance to start...")
	err := waitForZoneOperationState("DONE", c.Zone, instanceOperationName, client, c.stateTimeout)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	err = waitForInstanceState("RUNNING", c.Zone, instanceName, client, c.stateTimeout)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("instance_ip", "not yet")
	return multistep.ActionContinue
}

func (s *stepInstanceInfo) Cleanup(state multistep.StateBag) {
	// no cleanup
}
