// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package googlecompute

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepCreateSSHKey struct{}

func (s *stepCreateSSHKey) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	ui.Say("Creating temporary ssh key for instance...")
	priv, err := rsa.GenerateKey(rand.Reader, 2014)

	if err != nil {
		err := fmt.Errorf("Error creating temporary ssh key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	// ASN.1 DER encoded form
	priv_der := x509.MarshalPKCS1PrivateKey(priv)
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}
	// Marshal the public key into SSH compatible format
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		err := fmt.Errorf("Error creating temporary ssh key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("ssh_private_key", string(pem.EncodeToMemory(&priv_blk)))
	state.Put("ssh_public_key", string(ssh.MarshalAuthorizedKey(pub)))
	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	// no clean up, keys are created directly on the GCE instance.
}
