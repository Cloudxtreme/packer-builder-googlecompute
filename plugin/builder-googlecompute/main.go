package main

import (
	"github.com/mitchellh/packer/builder/googlecompute"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	plugin.ServeBuilder(new(googlecompute.Builder))
}
