package kable

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var LocalConceptDataSource = func() *schema.Resource {
	out := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path to the concept directory.",
			},
			"values": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "The values map to render the concept.",
			},
			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "yaml",
				Description: "The type of the rendered output. (defaults to 'yaml')",
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rendered output of the concept.",
			},
		},
		Read:        LocalConceptRead,
		Description: "The concept data source allows for rendering a concept from a repository",
	}

	return &out
}

var LocalConceptRead = func(d *schema.ResourceData, m interface{}) error {
	conceptPath := d.Get("path").(string)
	targetType := d.Get("target_type").(string)

	avs, err := assertValues(d)
	if err != nil {
		return err
	}

	if err := renderConcept(d, conceptPath, avs, targetType, true); err != nil {
		return err
	}

	return nil
}
