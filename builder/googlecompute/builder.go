// The googlecompute package contains a packer.Builder implementation that
// builds images for Google Compute Engine.
package googlecompute

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
)

// The unique ID for this builder
const BuilderId = "kelseyhightower.googlecompute"

type Builder struct {
	config config
	runner multistep.Runner
}

type config struct {
	common.PackerConfig `mapstructure:",squash"`
	ClientSecretsPath   string            `mapstructure:"client_secrets_path"`
	ImageName           string            `mapstructure:"image_name"`
	ImageDescription    string            `mapstructure:"image_description"`
	MachineType         string            `mapstructure:"machine_type"`
	Metadata            map[string]string `mapstructure:"metadata"`
	Network             string            `mapstructure:"network"`
	PrivateKeyPath      string            `mapstructure:"private_key_path"`
	ProjectId           string            `mapstructure:"project_id"`
	SourceImage         string            `mapstructure:"source_image"`
	SSHUsername         string            `mapstructure:"ssh_username"`
	SSHPort             uint              `mapstructure:"ssh_port"`
	RawSSHTimeout       string            `mapstructure:"ssh_timeout"`
	RawStateTimeout     string            `mapstructure:"state_timeout"`
	Tags                []string          `mapstructure:"tags"`
	Zone                string            `mapstructure:"zone"`
	// Private configuration settings not seen by the user.
	clientSecrets   *clientSecrets
	instanceName    string
	privateKeyBytes []byte
	sshTimeout      time.Duration
	stateTimeout    time.Duration
	tpl             *packer.ConfigTemplate
}

func (b *Builder) Prepare(raws ...interface{}) error {
	// Load the packer config.
	md, err := common.DecodeConfig(&b.config, raws...)
	if err != nil {
		return err
	}
	b.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return err
	}
	b.config.tpl.UserVars = b.config.PackerUserVars

	errs := common.CheckUnusedConfig(md)
	// Collect errors if any.
	if err := common.CheckUnusedConfig(md); err != nil {
		return err
	}
	// Set defaults.
	if b.config.Network == "" {
		b.config.Network = "default"
	}
	if b.config.ImageName == "" {
		// Default to packer-{{ unix timestamp (utc) }}
		b.config.ImageName = "packer-{{timestamp}}"
	}
	if b.config.MachineType == "" {
		b.config.MachineType = "n1-standard-1"
	}
	if b.config.RawSSHTimeout == "" {
		b.config.RawSSHTimeout = "5m"
	}
	if b.config.RawStateTimeout == "" {
		b.config.RawStateTimeout = "5m"
	}
	if b.config.SSHUsername == "" {
		b.config.SSHUsername = "root"
	}
	if b.config.SSHPort == 0 {
		b.config.SSHPort = 22
	}
	// Process required parameters.
	// client_secrets_path is required.
	if b.config.ClientSecretsPath == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a client_secrets_path must be specified"))
	}
	// private_key_path is required.
	if b.config.PrivateKeyPath == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a private_key_path must be specified"))
	}
	// project_id is required.
	if b.config.ProjectId == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a project_id must be specified"))
	}
	if b.config.SourceImage == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a source_image must be specified"))
	}
	if b.config.Zone == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a zone must be specified"))
	}
	// Process timeout settings.
	sshTimeout, err := time.ParseDuration(b.config.RawSSHTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing ssh_timeout: %s", err))
	}
	b.config.sshTimeout = sshTimeout
	// Set the state timeout.
	stateTimeout, err := time.ParseDuration(b.config.RawStateTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing state_timeout: %s", err))
	}
	b.config.stateTimeout = stateTimeout
	// Load the client secrets file.
	cs, err := loadClientSecrets(b.config.ClientSecretsPath)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing client secrets file: %s", err))
	}
	b.config.clientSecrets = cs
	// Load the private key.
	privateKeyBytes, err := ioutil.ReadFile(b.config.PrivateKeyPath)
	if err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Failed parsing client secrets file: %s", err))

	}
	b.config.privateKeyBytes = privateKeyBytes
	// Check for any errors.
	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Initialize the Google Compute Engine api.
	client, err := New(b.config.ProjectId, b.config.Zone, b.config.clientSecrets, b.config.privateKeyBytes)
	if err != nil {
		log.Println("Failed to create the Google Compute Engine client.")
		return nil, err
	}
	// Set up the state.
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)
	// Build the steps
	steps := []multistep.Step{
		new(stepCreateSSHKey),
		new(stepCreateInstance),
		new(stepInstanceInfo),
		&common.StepConnectSSH{
			SSHAddress:     sshAddress,
			SSHConfig:      sshConfig,
			SSHWaitTimeout: 5 * time.Minute,
		},
		new(common.StepProvision),
		new(stepCreateImage),
	}
	// Run the steps
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}
	b.runner.Run(state)
	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}
	if _, ok := state.GetOk("image_name"); !ok {
		log.Println("Failed to find image_name in state. Bug?")
		return nil, nil
	}
	artifact := &Artifact{
		imageName: state.Get("image_name").(string),
		client:    client,
	}
	return artifact, nil
}

func (b *Builder) Cancel() {
	return
}
