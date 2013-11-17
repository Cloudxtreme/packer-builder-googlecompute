# Google Compute Builder

Type: `googlecompute`

The `googlecompute` builder is able to create new images for use with
[Google Compute](https://cloud.google.com/products/compute-engine).

## Install

Download and build Packer from source as described [here](https://github.com/mitchellh/packer#developing-packer).

Next, clone this repository into `$GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute`.  Then build the packer-builder-googlecompute binary:

```
cd $GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute
go build -o /usr/local/packer/packer-builder-googlecompute \
plugin/builder-googlecompute/main.go
```

Now [configure Packer](http://www.packer.io/docs/other/core-configuration.html) to pick up the new builder:

`~/.packerconfig`

```
{
  "builders": {
    "googlecompute": "/usr/local/packer/packer-builder-googlecompute"
  }
}
```

## Basic Example

```JSON
{
  "builders": [{
    "type": "googlecompute",
    "bucket_name": "packer-images",
    "client_secrets_path": "/path/client_secrets.json",
    "private_key_path": "/path/private.key",
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
* `client_secrets_path` (string) - The client secrets file.
* `private_key_path` (string) - The Google Compute service account private key.
* `project_id` (string) - The Google Compute project id.
* `source_image` (string) - The source image to use. For example "debian-7-wheezy-v20131014".
* `zone` (string) - The Google Compute zone.

### Optional parameters:

* `image_name` (string) - The name of the resulting image that will appear in your account. This must be unique. To help make this unique, use a function like timestamp.
* `image_description` (string) - The description of the resulting image.
* `preferred_kernel` (string) - The preferred kernel to use with this image. Defaults to "gce-no-conn-track-v20130813".
* `machine_type` (string) - The machine type to use when building the image. The machine type must have a scratch disk. Defaults to "n1-standard-1".
* `network` (string) - The Google Compute network. Defaults to "default".
* `ssh_port` (int) - The port that SSH will be available on. Defaults to port 22.
* `ssh_timeout` (string) - The time to wait for SSH to become available before timing out. The format of this value is a duration such as "5s" or "5m". The default SSH timeout is "1m".
* `ssh_username` (string) - The username to use in order to communicate over SSH to the running instance. Default is "root".
* `state_timeout` (string) - The time to wait, as a duration string, for a instance to enter a desired state before timing out. The default state timeout is "6m".
