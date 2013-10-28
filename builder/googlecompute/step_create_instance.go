package googlecompute

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"
	"time"

	"code.google.com/p/google-api-go-client/compute/v1beta16"
)

type stepCreateInstance struct {
	instanceName string
}

func (s *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*GoogleComputeClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)

	ui.Say("Creating instance...")

	// Some random instance name as it's temporary
	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())

	// Build up the instance config. We need fully-qualified urls for the image and network.
	instanceConfig := &InstanceConfig{
		Description: "New instance created by Packer",
		Name:        name,
	}
	// Validate the zone.
	zone, err := client.GetZone(c.Zone)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Get the image
	image, err := client.GetImage(c.SourceImage)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	instanceConfig.Image = image.SelfLink

	machineType, err := client.GetMachineType(c.MachineType, zone.Name)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt

	}
	instanceConfig.MachineType = machineType.SelfLink

	network, err := client.GetNetwork(c.Network)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	networkInterface := NewNetworkInterface(network, true)
	networkInterfaces := []*compute.NetworkInterface{
		networkInterface,
	}
	instanceConfig.NetworkInterfaces = networkInterfaces

	// Create the instance based on configuration
	operation, err := client.CreateInstance(zone.Name, instanceConfig)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Wait for instance to go up.
	ui.Say("Waiting for the instance to start...")

	// Check the operation from the create instance call, then check the
	// instance status.
	for {
		status, err := client.ZoneOperationStatus(zone.Name, operation.Name)
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
		status, err := client.InstanceStatus(zone.Name, name)
		if err != nil {
			err := fmt.Errorf("Error creating instance: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		if status == "RUNNING" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	// We use this in cleanup
	s.instanceName = name
	// Store the image name for later
	state.Put("instance_name", name)
	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	if s.instanceName == "" {
		return
	}
	client := state.Get("client").(*GoogleComputeClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)

	// Destroy the instance we just created
	ui.Say("Destroying instance...")

	operation, err := client.DeleteInstance(c.Zone, s.instanceName)
	if err != nil {
		ui.Error(fmt.Sprintf("Error destroying instance. Please destroy it manually: %v", s.instanceName))
	}
	ui.Say("Waiting for the instance to be deleted...")
	for {
		status, err := client.ZoneOperationStatus(c.Zone, operation.Name)
		if err != nil {
			ui.Error(fmt.Sprintf("Error destroying instance. Please destroy it manually: %v", s.instanceName))
		}
		if status == "DONE" {
			break
		}
	}
	return
}
