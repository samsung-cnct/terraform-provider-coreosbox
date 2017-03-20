package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCoreOSBoxAmi() *schema.Resource {
	return &schema.Resource{
		Create: CreateAmi,
		Delete: DeleteAmi,
		Exists: ExistsAmi,
		Read:   ReadAmi,

		Schema: map[string]*schema.Schema{
			"channel": &schema.Schema{
				Type:        schema.TypeString,
				Description: "CoreOS update channel",
				Default:     "stable",
				Optional:    true,
				ForceNew:    true,
			},
			"virtualization": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Virtualization type. Applies to ami boxes: pv or hvm",
				Default:     "hvm",
				Optional:    true,
				ForceNew:    true,
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Region. Applies to ami boxes.",
				Default:     "us-west-2",
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

func CreateAmi(d *schema.ResourceData, meta interface{}) error {
	version, err := getVerInfo(d)
	if err != nil {
		return err
	}
	d.Set("version_out", version)

	box_info, err := getAmi(d)
	if err != nil {
		return err
	}
	d.Set("box_string", strings.TrimSpace(box_info))

	d.SetId(getID(d))
	return nil
}

func DeleteAmi(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func ExistsAmi(d *schema.ResourceData, meta interface{}) (bool, error) {
	return getID(d) == d.Id(), nil
}

func ReadAmi(d *schema.ResourceData, meta interface{}) error {
	version, err := getVerInfo(d)
	if err != nil {
		return err
	}
	d.Set("version_out", version)

	box_info, err := getAmi(d)
	if err != nil {
		return err
	}
	d.Set("box_string", strings.TrimSpace(box_info))

	d.SetId(getID(d))
	return nil
}

func getAmi(d *schema.ResourceData) (val string, err error) {
	type (
		ami struct {
			Name string `json:"name"`
			PV   string `json:"pv"`
			HVM  string `json:"hvm"`
		}

		amiInfo struct {
			AMIs []ami `json:"amis"`
		}
	)

	url := fmt.Sprintf(
		"http://%s.release.core-os.net/amd64-usr/%s/coreos_production_ami_all.json",
		d.Get("channel").(string),
		d.Get("version").(string))

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data amiInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	virt := d.Get("virtualization").(string)
	region := d.Get("region").(string)

	for _, a := range data.AMIs {
		if a.Name == region {
			switch virt {
			case "pv":
				return a.PV, nil
			case "hvm":
				return a.HVM, nil
			default:
				return "", errors.New("Unknown virtualization type")
			}
		}
	}
	return "", errors.New("No ami found")
}
