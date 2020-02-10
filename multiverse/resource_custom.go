package multiverse

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCustom() *schema.Resource {
	return &schema.Resource{
		Create: onCreate,
		Read:   onRead,
		Update: onUpdate,
		Delete: onDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"executor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"script": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"id_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"resource": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},

			"deep_object": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func onCreate(d *schema.ResourceData, m interface{}) error {
	return do("create", d, m)
}

func onRead(d *schema.ResourceData, m interface{}) error {
	return do("read", d, m)
}

func onUpdate(d *schema.ResourceData, m interface{}) error {
	return do("update", d, m)
}

func onDelete(d *schema.ResourceData, m interface{}) error {
	return do("delete", d, m)
}

// RPCCommand - standart JSON struct for passing data in STDIN
type RPCCommand struct {
	ID         string `json:"ID,omitempty"`
	Payload    string `json:"Payload,omitempty"`
	DeepObject string `json:"DeepObject,omitempty"`
}

func do(event string, d *schema.ResourceData, m interface{}) error {
	log.Printf("Executing: %s %s %s %s", d.Get("executor"), d.Get("script"), event, d.Get("data"))

	cmd := exec.Command(d.Get("executor").(string), d.Get("script").(string), event)

	rpcCmd := RPCCommand{
		ID:      d.Id(),
		Payload: d.Get("data").(string),
	}

	jsonStr, err := json.Marshal(&rpcCmd)
	if err != nil {
		return err
	}
	cmd.Stdin = bytes.NewReader([]byte(jsonStr))
	result, err := cmd.Output()

	if err != nil {
		return err
	}

	if err == nil {
		ret := RPCCommand{}
		err := json.Unmarshal([]byte(result), &ret)
		if err != nil {
			return err
		}
		if event == "delete" {
			d.SetId("")
		} else {
			d.SetId(ret.ID)
			d.Set("deep_object", ret.DeepObject)
			var resource map[string]interface{}
			json.Unmarshal([]byte(result), &resource)
			d.Set("resource", resource)
		}
	}

	// var resource map[string]interface{}
	// err = json.Unmarshal([]byte(result), &resource)
	// if err == nil {
	// 	if event == "delete" {
	// 		d.SetId("")
	// 	} else {
	// 		key := d.Get("id_key").(string)
	// 		d.Set("resource", resource)
	// 		d.SetId(resource[key].(string))
	// 	}
	// }

	return err
}
