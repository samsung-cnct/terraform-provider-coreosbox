# CoreOS Version provider for Terraform

This [Terraform](http://terraform.io) provider is for dynamically finding version number and box information for a given IaaS, update channel and version (or current version)

## Status

Development/Testing

Currently only supports ami, gce and vagrant box info

## Install - Build

```
go get
go build
go install
```

## With Homebrew

```
$ brew tap 'samsung-cnct/terraform-provider-coreosbox'
$ brew install terraform-provider-coreosbox
```

## Usage


```
resource "coreosbox_ami" "current_ami" {
  channel = "alpha"
  virtualization = "hvm"
  region = "us-east-1"
  version = "current"
}

resource "coreosbox_ami" "versioned_ami" {
  channel = "stable"
  virtualization = "hvm"
  region = "us-east-1"
  version = "647.2.0"
}

resource "coreosbox_gce" "current_gce" {
  channel = "alpha"
  version = "current"
}

resource "coreosbox_gce" "versioned_gce" {
  channel = "stable"
  version = "647.2.0"
}

resource "coreosbox_vagrant" "current_vagrant" {
  channel = "alpha"
  version = "current"
  hypervisor = "virtualbox"
}

resource "coreosbox_vagrant" "versioned_vagrant" {
  channel = "stable"
  version = "647.2.0"
  hypervisor = "virtualbox"
}

resource "coreosbox_vagrant" "current_vagrant_vmware" {
  channel = "alpha"
  version = "current"
  hypervisor = "vmware"
}

resource "coreosbox_vagrant" "versioned_vagrant_vmware" {
  channel = "stable"
  version = "647.2.0"
  hypervisor = "vmware"
}

output "info_ami" {
    value = "Version: ${coreosbox_ami.current_ami.version_out}, ami: ${coreosbox_ami.current_ami.box_string}." 
}

output "info_ami_versioned" {
    value = "Version: ${coreosbox_ami.versioned_ami.version_out}, ami: ${coreosbox_ami.versioned_ami.box_string}." 
}

output "info_gce" {
    value = "Version: ${coreosbox_gce.current_gce.version_out}, box: ${coreosbox_gce.current_gce.box_string}." 
}

output "info_gce_versioned" {
    value = "Version: ${coreosbox_gce.versioned_gce.version_out}, box: ${coreosbox_gce.versioned_gce.box_string}." 
}

output "info_vagrant" {
    value = "Version: ${coreosbox_vagrant.current_vagrant.version_out}, box: ${coreosbox_vagrant.current_vagrant.box_string}." 
}

output "info_vagrant_versioned" {
    value = "Version: ${coreosbox_vagrant.versioned_vagrant.version_out}, box: ${coreosbox_vagrant.versioned_vagrant.box_string}." 
}

output "info_vmware" {
    value = "Version: ${coreosbox_vagrant.current_vagrant_vmware.version_out}, box: ${coreosbox_vagrant.current_vagrant_vmware.box_string}." 
}

output "info_vmware_versioned" {
    value = "Version: ${coreosbox_vagrant.versioned_vagrant_vmware.version_out}, box: ${coreosbox_vagrant.versioned_vagrant_vmware.box_string}." 
}
```

## Cutting release
This is a manual release process, we may automate it in the future if there is a need.

Steps:
1. build linux executable:
`GOOS=linux GOARCH=amd64 go build`
2. tar linux executable
`tar -cf terraform-provider-coreosbox_linux_amd64.tar terraform-provider-coreosbox`
3. gzip linux executable
`gzip terraform-provider-coreosbox_linux_amd64.tar`
4. repeat above for darwin build
5. click 'Draft a new release' on releases page
6. fill out form
7. Add both tar.gz files