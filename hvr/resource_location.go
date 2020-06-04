package hvr

import (
	"fmt"
	"regexp"

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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hvrLocationNameValidator(),
			},
			"class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"db2i", "postgresql"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"database_connection": {
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
			"remote_machine": {
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
	service := meta.(Service)

	location := expandLocation(d)

	if err := service.NewLocation(location); err != nil {
		return fmt.Errorf("error creating HVR location %+v: %s", d, err)
	}

	d.SetId(location.Name)
	return resourceHVRLocationRead(d, meta)
}

func resourceHVRLocationRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	location, err := service.GetLocation(d.Id())

	if err != nil {
		if err.Error() == "location not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading HVR location %s: %s", d.Id(), err)
	}

	d.Set("name", location.Name)
	d.Set("class", location.Class)
	d.Set("description", location.Description)
	d.Set("database_connection", flattenDatabaseConnection(location))
	d.Set("remote_machine", flattenRemoteMachine(location))

	return nil
}

func resourceHVRLocationUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	location := expandLocation(d)

	if err := service.UpdateLocation(location); err != nil {
		return fmt.Errorf("error updating HVR location %+v: %s", d, err)
	}

	return resourceHVRLocationRead(d, meta)
}

func resourceHVRLocationDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	if err := service.DeleteLocation(d.Id()); err != nil {
		return fmt.Errorf("error deleting HVR location %+v: %s", d, err)
	}

	d.SetId("")
	return nil
}

func hvrLocationNameValidator() schema.SchemaValidateFunc {
	r := regexp.MustCompile("^[[:alpha:]][[:alnum:]_]*$")
	violationMessage := "names must be composed of alphabetics followed by alphanumerics or underscores"

	return validation.All(
		validation.StringLenBetween(1, 5),
		validation.StringMatch(r, violationMessage),
	)
}
