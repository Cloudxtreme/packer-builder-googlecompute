// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package googlecompute

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"

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
	name := fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	// Build up the instance config.
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
	// Set the source image. Must be a fully-qualified URL.
	image, err := client.GetImage(c.SourceImage)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	instanceConfig.Image = image.SelfLink
	// Set the machineType. Must be a fully-qualified URL.
	machineType, err := client.GetMachineType(c.MachineType, zone.Name)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt

	}
	instanceConfig.MachineType = machineType.SelfLink
	// Set up the Network Interface.
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
	// Add the metadata, which also setups up the ssh key.
	metadata := make(map[string]string)
	sshPublicKey := state.Get("ssh_public_key").(string)
	metadata["sshKeys"] = fmt.Sprintf("%s:%s", c.SSHUsername, sshPublicKey)
	instanceConfig.Metadata = MapToMetadata(metadata)
	// Add the default service so we can create an image of the machine and
	// upload it to cloud storage.
	defaultServiceAccount := NewServiceAccount("default")
	serviceAccounts := []*compute.ServiceAccount{
		defaultServiceAccount,
	}
	instanceConfig.ServiceAccounts = serviceAccounts
	// Create the instance based on configuration
	operation, err := client.CreateInstance(zone.Name, instanceConfig)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("Waiting for the instance to be created...")
	err = waitForZoneOperationState("DONE", c.Zone, operation.Name, client, c.stateTimeout)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	// Update the state.
	state.Put("instance_name", name)
	s.instanceName = name
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
