// The googlecompute package contains a packer.Builder implementation that
// builds images for Google Compute Engine.
package googlecompute

import (
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"log"
	"time"
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
