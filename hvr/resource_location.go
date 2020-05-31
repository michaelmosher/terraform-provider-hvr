package hvr

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceHVRLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceHVRLocationCreate,
		Read:   resourceHVRLocationRead,
		Update: resourceHVRLocationUpdate,
		Delete: resourceHVRLocationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 5),
			},
			"class": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"db2i", "postgresql"}, false),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"database_connection": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IsPortNumber,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"remote_machine": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceHVRLocationCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)
	return resourceHVRLocationRead(d, meta)
}

func resourceHVRLocationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceHVRLocationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceHVRLocationRead(d, meta)
}

func resourceHVRLocationDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
