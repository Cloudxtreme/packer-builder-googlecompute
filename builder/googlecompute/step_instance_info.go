// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

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
	err := waitForInstanceState("RUNNING", c.Zone, instanceName, client, c.stateTimeout)
	if err != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ip, err := client.GetNatIP(c.Zone, instanceName)
	if err != nil {
		err := fmt.Errorf("Error retrieving instance nat ip address: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("instance_ip", ip)
	return multistep.ActionContinue
}

func (s *stepInstanceInfo) Cleanup(state multistep.StateBag) {
	// no cleanup
}
