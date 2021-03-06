package hvr

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/michaelmosher/go-hvr-sdk/hvrhub"
)

func resourceHVRChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceHVRChannelCreate,
		Read:   resourceHVRChannelRead,
		Update: resourceHVRChannelUpdate,
		Delete: resourceHVRChannelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hvrChannelNameValidator(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceHVRChannelCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channel := hvrhub.Channel{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	err := service.NewChannel(channel)

	if err != nil {
		return fmt.Errorf("error creating HVR channel %+v: %s", d, err)
	}

	d.SetId(channel.Name)
	return resourceHVRChannelRead(d, meta)
}

func resourceHVRChannelRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channel, err := service.GetChannel(d.Id())

	if err != nil {
		if err.Error() == "channel not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading HVR channel %s: %s", d.Id(), err)
	}

	d.Set("name", channel.Name)
	d.Set("description", channel.Description)

	return nil
}

func resourceHVRChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channel := hvrhub.Channel{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if err := service.UpdateChannel(channel); err != nil {
		return fmt.Errorf("error updating HVR channel %+v: %s", d, err)
	}

	return resourceHVRChannelRead(d, meta)
}

func resourceHVRChannelDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	if err := service.DeleteChannel(d.Id()); err != nil {
		return fmt.Errorf("error deleting HVR channel %+v: %s", d, err)
	}

	d.SetId("")
	return nil
}

func hvrChannelNameValidator() schema.SchemaValidateFunc {
	r := regexp.MustCompile("^[[:alpha:]][[:alnum:]_]*$")
	violationMessage := "names must be composed of alphabetics followed by alphanumerics or underscores"

	return validation.All(
		validation.StringLenBetween(1, 12),
		validation.StringMatch(r, violationMessage),
	)
}
