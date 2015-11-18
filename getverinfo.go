package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "strings"

  "github.com/dlintw/goconf"
  "github.com/hashicorp/terraform/helper/schema"
)

func getVerInfo(d *schema.ResourceData) (val string, err error) {
  url := fmt.Sprintf(
    "http://%s.release.core-os.net/amd64-usr/%s/version.txt", 
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
  verInfo, err := goconf.ReadConfigBytes(bodyBytes)
  if err != nil {
    return "", err
  }

  return verInfo.GetString("default", "COREOS_VERSION")
}

func getID(d *schema.ResourceData) string {
  channel := d.Get("channel").(string)
  v := d.Get("version").(string)

  return strings.Join([]string{channel, v}, ":")
}

