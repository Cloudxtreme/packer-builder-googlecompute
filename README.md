# Google Compute Builder

Type: `googlecompute`

### Not ready for use quite yet, 80% complete.

The `googlecompute` build is able to create new images for use with
[Google Compute](https://cloud.google.com/products/compute-engine).

## Install

Download and build Packer from source as described [here](https://github.com/mitchellh/packer#developing-packer).

Next, clone this repository into `$GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute`.  Then build the packer-builder-googlecompute binary:

```
cd $GOPATH/src/github.com/kelseyhightower/packer-builder-googlecompute
go build -o /usr/local/packer/packer-builder-googlecompute \
plugin/packer-builder-googlecompute/main.go
```

Now [configure Packer](http://www.packer.io/docs/other/core-configuration.html) to pick up the new builder:

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
  "type": "googlecompute",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "machine_type": "n1-standard-1-d",
  "client_email": "XXXXXXXXXXXXXXX@developer.gserviceaccount.com",
  "client_id": "XXXXXXXXXXXXXXX.apps.googleusercontent.com",
  "private_key_path": "/path/to/service_account/privatekey_pem_file",
  "source_image": "debian-7-wheezy-v20130926",
  "ssh_username": "",
  "ssh_timeout": "5m"
}
```
