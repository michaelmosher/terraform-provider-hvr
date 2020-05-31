package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/michaelmosher/terraform-provider-hvr/hvr"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hvr.Provider,
	})
}
