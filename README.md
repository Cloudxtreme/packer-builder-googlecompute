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

Required parameters:

* `bucket_name` (string) - The Google Cloud Storage bucket to store machine images.
* `client_secrets_path` (string) - The Google Compute client secrets file.
* `private_key_path` (string) - The Google Compute service account private key.
* `project_id` (string) - The Google Compute project id.
* `source_image` (string) - The source image to use. For example debian-7-wheezy-v20131014.
* `zone` (string) - The Google Compute zone.
