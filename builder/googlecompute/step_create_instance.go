package googlecompute

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"
)

type stepCreateInstance struct {
	instanceName uint
}

func (s *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*GoogleComputeClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)
	sshPublicKey := state.Get("ssh_public_key").(string)

	ui.Say("Creating instance...")

	// Some random instance name as it's temporary
	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())

	// Build up the instance config. We need fully-qualified urls for the image and network.
	instanceConfig := &InstanceConfig{
		Description: "New instance created by Packer",
		Name:        name,
	}
	// Validate the zone.
	zone, err := g.GetZone(zoneName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get the image
	image, err := g.GetImage(imageName)
	if err != nil {
		log.Fatal(err.Error())
	}
	instanceConfig.Image = image.SelfLink

	machineType, err := g.GetMachineType(machineTypeName, zone.Name)
	if err != nil {
		log.Fatal(err.Error())
	}
	instanceConfig.MachineType = machineType.SelfLink

	network, err := g.GetNetwork(networkName)
	if err != nil {
		log.Fatal(err.Error())
	}
	networkInterface := googlecompute.NewNetworkInterface(network, true)
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
	s.imageName = name
	// Store the image name for later
	state.Put("image_name", imageName)
	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	if s.imageName == "" {
		return
	}
	client := state.Get("client").(*GoogleComputeClient)
	ui := state.Get("ui").(packer.Ui)
	c := state.Get("config").(config)

	// Destroy the instance we just created
	ui.Say("Destroying instance...")

	operation, err = client.DeleteInstance(zone.Name, s.imageName)
	if err != nil {
		ui.Error(fmt.Sprintf("Error destroying instance. Please destroy it manually: %v", s.imageName))
	}
	ui.Say("Waiting for the instance to be deleted...")
	for {
		status, err := g.ZoneOperationStatus(zone.Name, operation.Name)
		if err != nil {
			ui.Error(fmt.Sprintf("Error destroying instance. Please destroy it manually: %v", s.imageName))
		}
		if status == "DONE" {
			break
		}
	}
	return
}
