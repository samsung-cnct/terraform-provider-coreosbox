package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"coreosbox_ami":     resourceCoreOSBoxAmi(),
			"coreosbox_gce":     resourceCoreOSBoxGce(),
			"coreosbox_vagrant": resourceCoreOSBoxVagrant(),
		},
	}
}
