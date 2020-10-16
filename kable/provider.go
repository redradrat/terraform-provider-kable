package kable

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var Provider = func() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"kable_concept":       ConceptDataSource(),
			"kable_local_concept": LocalConceptDataSource(),
		},
	}
}
