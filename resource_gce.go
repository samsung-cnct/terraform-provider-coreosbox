package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"

  "github.com/hashicorp/terraform/helper/schema"
)


func resourceCoreOSBoxGce() *schema.Resource {
  return &schema.Resource{
    Create: CreateGce,
    Delete: DeleteGce,
    Exists: ExistsGce,
    Read:   ReadGce,

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

func CreateGce(d *schema.ResourceData, meta interface{}) error {
  version, err := getVerInfo(d)
  if err != nil {
    return err
  }
  d.Set("version_out", version)

  box_info, err := getGce(d)
  if err != nil {
    return err
  }
  d.Set("box_string", strings.TrimSpace(box_info))

  d.SetId(getID(d))
  return nil
}

func DeleteGce(d *schema.ResourceData, meta interface{}) error {
  d.SetId("")
  return nil
}

func ExistsGce(d *schema.ResourceData, meta interface{}) (bool, error) {
  return getID(d) == d.Id(), nil
}

func ReadGce(d *schema.ResourceData, meta interface{}) error {
  version, err := getVerInfo(d)
  if err != nil {
    return err
  }
  d.Set("version_out", version)

  box_info, err := getGce(d)
  if err != nil {
    return err
  }
  d.Set("box_string", strings.TrimSpace(box_info))

  d.SetId(getID(d))
  return nil
}

func getGce(d *schema.ResourceData) (val string, err error) {
  url := fmt.Sprintf(
    "http://%s.release.core-os.net/amd64-usr/%s/coreos_production_gce.txt", 
    d.Get("channel").(string), 
    d.Get("version").(string))

  resp, err := http.Get(url)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()

  bodyBytes, err := ioutil.ReadAll(resp.Body) 
  if err != nil {
    return "", err
  }

  return string(bodyBytes[:]), nil
}