// The googlecompute package contains a packer.Builder implementation that
// builds images for Google Compute Engine.
package googlecompute

import (
	"errors"
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"log"
	"os"
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
	tpl                 *packer.ConfigTemplate
}

func (b *Builder) Prepare(raws ...interface{}) error {
	// Nothing yet.
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Nothing yet.
}

func (b *Builder) Cancel() {
	// Nothing yet.
}
