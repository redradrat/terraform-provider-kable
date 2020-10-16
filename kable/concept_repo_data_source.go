package kable

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

var ConceptRepoDataSource = func() *schema.Resource {
	out := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id for the repository",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The url for the repository",
			},
			"ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A git ref can be given to reference a specific version",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Should the repository required authentication, a username can be provided",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Should the repository required authentication, a password can be provided",
				Sensitive:   true,
			},
		},
		Description: "The concept repo data source allows referencing a concept repository",
	}

	return &out
}
