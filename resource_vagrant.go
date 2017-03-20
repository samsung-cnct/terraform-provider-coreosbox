package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCoreOSBoxVagrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateVagrant,
		Delete: DeleteVagrant,
		Exists: ExistsVagrant,
		Read:   ReadVagrant,

		Schema: map[string]*schema.Schema{
			"channel": &schema.Schema{
				Type:        schema.TypeString,
				Description: "CoreOS update channel",
				Default:     "stable",
				Optional:    true,
				ForceNew:    true,
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Description: "CoreOS version: current or version number",
				Default:     "current",
				Optional:    true,
				ForceNew:    true,
			},
			"hypervisor": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Vagrant provider: virtualbox or vmware",
				Default:     "virtualbox",
				Optional:    true,
				ForceNew:    true,
			},
			"version_out": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version string read from CoreOS releases",
			},
			"box_string": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Box information string. Ami id, gce path, vagrant box url",
			},
		},
	}
}

func CreateVagrant(d *schema.ResourceData, meta interface{}) error {
	version, err := getVerInfo(d)
	if err != nil {
		return err
	}
	d.Set("version_out", version)

	box_info, err := getVagrant(d)
	if err != nil {
		return err
	}
	d.Set("box_string", strings.TrimSpace(box_info))

	d.SetId(getID(d))
	return nil
}

func DeleteVagrant(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func ExistsVagrant(d *schema.ResourceData, meta interface{}) (bool, error) {
	return getID(d) == d.Id(), nil
}

func ReadVagrant(d *schema.ResourceData, meta interface{}) error {
	version, err := getVerInfo(d)
	if err != nil {
		return err
	}
	d.Set("version_out", version)

	box_info, err := getVagrant(d)
	if err != nil {
		return err
	}
	d.Set("box_string", strings.TrimSpace(box_info))

	d.SetId(getID(d))
	return nil
}

func getVagrant(d *schema.ResourceData) (val string, err error) {
	switch d.Get("hypervisor").(string) {
	case "virtualbox":
		return getVirtualboxVagrant(d)
	case "vmware":
		return getVmwareVagrant(d)
	}

	return "", errors.New("Unknown vagrant hypervisor")
}

func getVirtualboxVagrant(d *schema.ResourceData) (val string, err error) {
	url := fmt.Sprintf(
		"http://%s.release.core-os.net/amd64-usr/%s/coreos_production_vagrant.json",
		d.Get("channel").(string),
		d.Get("version").(string))

	return url, nil
}

func getVmwareVagrant(d *schema.ResourceData) (val string, err error) {
	url := fmt.Sprintf(
		"http://%s.release.core-os.net/amd64-usr/%s/coreos_production_vagrant_vmware_fusion.json.json",
		d.Get("channel").(string),
		d.Get("version").(string))

	return url, nil
}
