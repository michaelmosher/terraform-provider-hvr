package hvr

import (
	"encoding/base64"
	"fmt"

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
				ValidateFunc: validation.StringLenBetween(1, 12),
			},
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 11),
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
	err := service.NewLocationGroup(group)

	if err != nil {
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
	err := service.UpdateLocationGroup(group)

	if err != nil {
		return fmt.Errorf("error updating HVR location group %+v: %s", d, err)
	}

	return resourceHVRLocationGroupRead(d, meta)
}

func resourceHVRLocationGroupDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channelName := d.Get("channel_name").(string)
	groupName := d.Get("group_name").(string)

	err := service.DeleteLocationGroup(channelName, groupName)

	if err != nil {
		return fmt.Errorf("error deleting HVR location group %+v: %s", d, err)
	}

	d.SetId("")
	return nil
}
