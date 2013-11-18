# Google Compute Builder

[![Build Status](https://travis-ci.org/kelseyhightower/packer-builder-googlecompute.png?branch=master)](https://travis-ci.org/kelseyhightower/packer-builder-googlecompute)

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

### GCE Credentials

The `googlecompute` builder requires a GCE [service account](https://developers.google.com/console/help/#service_accounts). 

The client_secrets.json and privatekey.p12 are required:

* client_secret_XXXXXX-XXXXXX.apps.googleusercontent.com.json
* XXXXXX-privatekey.p12

The `XXXXXX-privatekey.p12` must be converted to pem format. This can
be done using the openssl commandline tool:

```Bash
openssl pkcs12 -in XXXXXX-privatekey.p12 -out XXXXXX-privatekey.pem
```

When prompted for "Enter Import Password", enter `notasecret`.

## Basic Example

```JSON
{
  "builders": [{
    "type": "googlecompute",
    "bucket_name": "packer-images",
    "client_secrets_file": "client_secret_XXXXXX-XXXXXX.apps.googleusercontent.com.json",
    "private_key_file": "XXXXXX-privatekey.pem",
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
* `passphrase` (string) - The passphrase to use if the `private_key_file` is encrypted.
* `preferred_kernel` (string) - The preferred kernel. Defaults to `gce-no-conn-track-v20130813`.
* `ssh_port` (int) - The SSH port. Defaults to `22`.
* `ssh_timeout` (string) - The time to wait for SSH to become available. Defaults to `1m`.
* `ssh_username` (string) - The SSH username. Defaults to `root`.
* `state_timeout` (string) - The time to wait for instance state changes. Defaults to `5m`.

> The machine type must have a scratch disk.

## Building

Clone this repository into `$GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute`.  Then build the packer-builder-googlecompute binary:

```
cd $GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute
go get
go build
```

Copy the results to the Packer install directory.

```
cp packer-builder-googlecompute /usr/local/packer/packer-builder-googlecompute
```
