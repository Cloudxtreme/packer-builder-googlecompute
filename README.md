# Google Compute Builder

Type: `googlecompute`

The `googlecompute` build is able to create new images for use with
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
