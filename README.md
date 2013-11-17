# Google Compute Builder

Type: `googlecompute`

The `googlecompute` builder is able to create new [images](https://developers.google.com/compute/docs/images)
for use with [Google Compute Engine](https://cloud.google.com/products/compute-engine).

## Install

Download a binary release from [Github](https://github.com/kelseyhightower/packer-builder-googlecompute/releases).
Extract then copy the `packer-builder-googlecompute` binary to the Packer installation directory.

```Bash
unzip packer-builder-googlecompute_0.1.0-beta2_darwin_amd64.zip
cp packer-builder-googlecompute /usr/local/packer/
```

> Packer version v0.3.11+ required.

## Configure

Enable the googlecompute builder in `~/.packerconfig`

```
{
  "builders": {
    "googlecompute": "/usr/local/packer/packer-builder-googlecompute"
  }
}
```

> See [configure Packer](http://www.packer.io/docs/other/core-configuration.html) for more info.

## Basic Example

```JSON
{
  "builders": [{
    "type": "googlecompute",
    "bucket_name": "packer-images",
    "client_secrets_file": "/path/client_secrets.json",
    "private_key_file": "/path/private.key",
    "project_id": "my-project",
    "source_image": "debian-7-wheezy-v20131014",
    "zone": "us-central1-a"
  }]
}
```

## Configuration Reference

The reference of available configuration options is listed below.

### Required parameters:

* `bucket_name` (string) - The Google Cloud Storage bucket to store images.
* `client_secrets_file` (string) - The client secrets file.
* `private_key_file` (string) - The service account private key.
* `project_id` (string) - The GCE project id.
* `source_image` (string) - The source image. Example `debian-7-wheezy-v20131014`.
* `zone` (string) - The GCE zone.

### Optional parameters:

* `image_name` (string) - The unique name of the resulting image. Defaults to `packer-{{timestamp}}`.
* `image_description` (string) - The description of the resulting image.
* `machine_type` (string) - The machine type. Defaults to `n1-standard-1`.
* `network` (string) - The Google Compute network. Defaults to `default`.
* `preferred_kernel` (string) - The preferred kernel. Defaults to `gce-no-conn-track-v20130813`.
* `ssh_port` (int) - The SSH port. Defaults to `22`.
* `ssh_timeout` (string) - The time to wait for SSH to become available. Defaults to `1m`.
* `ssh_username` (string) - The SSH username. Defaults to `root`.
* `state_timeout` (string) - The time to wait for instance state changes. Defaults to `5m`.

> The machine type must have a scratch disk.

## Building

Download and build Packer from source as described [here](https://github.com/mitchellh/packer#developing-packer).

Next, clone this repository into `$GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute`.  Then build the packer-builder-googlecompute binary:

```
cd $GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute
go build -o /usr/local/packer/packer-builder-googlecompute \
plugin/builder-googlecompute/main.go
```
