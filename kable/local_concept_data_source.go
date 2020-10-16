package kable

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var LocalConceptDataSource = func() *schema.Resource {
	return &schema.Resource{
		ReadContext: LocalConceptRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path to the concept directory.",
			},
			"inputs": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The key/value list of inputs to use.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "yaml",
				Description: "The type of the rendered output. (defaults to 'yaml')",
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,

				Description: "The rendered output of the concept.",
			},
		},
		Description: "The concept data source allows for rendering a concept from a repository",
	}
}

var LocalConceptRead = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conceptPath := d.Get("path").(string)
	targetType := d.Get("target_type").(string)

	avs, err := assertValues(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := renderConcept(d, conceptPath, avs, targetType, true); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
