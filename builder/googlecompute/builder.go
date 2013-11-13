// The googlecompute package contains a packer.Builder implementation that
// builds images for Google Compute Engine.
package googlecompute

import (
	"errors"
	"fmt"
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
	AuthURI             string            `mapstructure:"auth_uri"`
	ClientEmail         string            `mapstructure:"client_email"`
	ClientId            string            `mapstructure:"client_id"`
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
	TokenURI            string            `mapstructure:"token_uri"`
	Zone                string            `mapstructure:"zone"`
	sshTimeout          time.Duration
	stateTimeout        time.Duration
	clientSecrets       *clientSecrets
	privateKeyBytes     []byte
	tpl                 *packer.ConfigTemplate
	instanceName        string
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
	// Do we need to collect anything from the environment?

	if b.config.ImageName == "" {
		// Default to packer-{{ unix timestamp (utc) }}
		b.config.ImageName = "packer-{{timestamp}}"
	}
	if b.config.MachineType == "" {
		b.config.MachineType = "default image type"
	}
	if b.config.SSHUsername == "" {
		b.config.SSHUsername = "root"
	}
	// Default Instance Size?
	// Set the default SSH port
	if b.config.SSHPort == 0 {
		b.config.SSHPort = 22
	}

	// Still need to process user vars and template strings.

	// Required configurations that will display errors if not set.
	if b.config.AuthURI == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a auth_uri must be specified"))
	}
	if b.config.ClientEmail == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a client_email must be specified"))
	}
	if b.config.ClientId == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a client_id must be specified"))
	}
	if b.config.PrivateKeyPath == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a private_key_path must be specified"))
	}
	if b.config.ProjectId == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a project_id must be specified"))
	}
	if b.config.RawSSHTimeout == "" {
		b.config.RawSSHTimeout = "1m"
	}
	if b.config.RawStateTimeout == "" {
		b.config.RawStateTimeout = "6m"
	}
	if b.config.SourceImage == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a source_image must be specified"))
	}
	if b.config.TokenURI == "" {
		errs = packer.MultiErrorAppend(
			errs, errors.New("a token_uri must be specified"))
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
