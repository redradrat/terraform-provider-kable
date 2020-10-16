package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/redradrat/terraform-provider-kable/kable"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kable.Provider})
}
