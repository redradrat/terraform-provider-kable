package kable

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var Provider = func() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"kable_concept":       ConceptDataSource(),
			"kable_local_concept": LocalConceptDataSource(),
			"kable_concept_repo":  ConceptRepoDataSource(),
		},
	}
}
