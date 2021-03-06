package hvr

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/michaelmosher/go-hvr-sdk/hvrhub"
)

func resourceHVRLocationGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceHVRLocationGroupCreate,
		Read:   resourceHVRLocationGroupRead,
		Update: resourceHVRLocationGroupUpdate,
		Delete: resourceHVRLocationGroupDelete,

		Schema: map[string]*schema.Schema{
			"channel_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hvrChannelNameValidator(),
			},
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hvrGroupNameValidator(),
			},
			"group_description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceHVRLocationGroupCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	group := hvrhub.LocationGroup{
		ChannelName: d.Get("channel_name").(string),
		Name:        d.Get("group_name").(string),
		Description: d.Get("group_description").(string),
	}

	if err := service.NewLocationGroup(group); err != nil {
		return fmt.Errorf("error creating HVR location group %+v: %s", d, err)
	}

	d.SetId(base64.StdEncoding.EncodeToString([]byte(group.ChannelName + group.Name)))
	return resourceHVRLocationGroupRead(d, meta)
}

func resourceHVRLocationGroupRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channelName := d.Get("channel_name").(string)
	groupName := d.Get("group_name").(string)

	group, err := service.GetLocationGroup(channelName, groupName)

	if err != nil {
		if err.Error() == "location group not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading HVR location group %s: %s", d.Id(), err)
	}

	d.Set("channel_name", group.ChannelName)
	d.Set("group_name", group.Name)
	d.Set("group_description", group.Description)

	return nil
}

func resourceHVRLocationGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	group := hvrhub.LocationGroup{
		ChannelName: d.Get("channel_name").(string),
		Name:        d.Get("group_name").(string),
		Description: d.Get("group_description").(string),
	}

	if err := service.UpdateLocationGroup(group); err != nil {
		return fmt.Errorf("error updating HVR location group %+v: %s", d, err)
	}

	return resourceHVRLocationGroupRead(d, meta)
}

func resourceHVRLocationGroupDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channelName := d.Get("channel_name").(string)
	groupName := d.Get("group_name").(string)

	if err := service.DeleteLocationGroup(channelName, groupName); err != nil {
		return fmt.Errorf("error deleting HVR location group %+v: %s", d, err)
	}

	d.SetId("")
	return nil
}

func hvrGroupNameValidator() schema.SchemaValidateFunc {
	r := regexp.MustCompile("^[[:upper:]0-9_]+$")
	violationMessage := "only uppercase letters, underscores and numbers are allowed."

	return validation.All(
		validation.StringLenBetween(1, 11),
		validation.StringMatch(r, violationMessage),
	)
}
