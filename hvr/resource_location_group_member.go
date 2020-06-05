package hvr

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/michaelmosher/go-hvr-sdk/hvrhub"
)

func resourceHVRLocationGroupMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceHVRLocationGroupMemberCreate,
		Read:   resourceHVRLocationGroupMemberRead,
		Delete: resourceHVRLocationGroupMemberDelete,

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
			"location_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: hvrLocationNameValidator(),
			},
		},
	}
}

func resourceHVRLocationGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	member := hvrhub.LocationGroupMember{
		ChannelName:       d.Get("channel_name").(string),
		LocationGroupName: d.Get("group_name").(string),
		LocationName:      d.Get("location_name").(string),
	}

	if err := service.NewLocationGroupMember(member); err != nil {
		return fmt.Errorf("error creating HVR location group %+v: %s", d, err)
	}

	d.SetId(
		base64.StdEncoding.EncodeToString(
			[]byte(member.ChannelName + member.LocationGroupName + member.LocationName),
		),
	)
	return resourceHVRLocationGroupRead(d, meta)
}

func resourceHVRLocationGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channelName := d.Get("channel_name").(string)
	groupName := d.Get("group_name").(string)
	locationName := d.Get("location_name").(string)

	// Just making sure the record still exists.
	_, err := service.GetLocationGroupMember(channelName, groupName, locationName)

	if err != nil {
		if err.Error() == "location group member not found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading HVR location group member %s: %s", d.Id(), err)
	}

	return nil
}

func resourceHVRLocationGroupMemberDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(Service)

	channelName := d.Get("channel_name").(string)
	groupName := d.Get("group_name").(string)
	locationName := d.Get("location_name").(string)

	if err := service.DeleteLocationGroupMember(channelName, groupName, locationName); err != nil {
		return fmt.Errorf("error deleting HVR location group member %+v: %s", d, err)
	}

	d.SetId("")
	return nil
}
