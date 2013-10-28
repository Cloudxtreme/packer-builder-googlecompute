package googlecompute

import (
	"fmt"
	"time"
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
	// Check the operation from the create instance call, then check the
	// instance status.
	for {
		status, err := client.ZoneOperationStatus(c.Zone, instanceOperationName)
		if err != nil {
			err := fmt.Errorf("Error creating instance: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		if status == "DONE" {
			break
		}
	}
	for {
		status, ip, err := client.InstanceStatus(c.Zone, instanceName)
		if err != nil {
			err := fmt.Errorf("Error creating instance: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		if status == "RUNNING" {
			state.Put("instance_ip", ip)
			break
		}
		time.Sleep(10 * time.Second)
	}
	return multistep.ActionContinue
}

func (s *stepInstanceInfo) Cleanup(state multistep.StateBag) {
	// no cleanup
}
