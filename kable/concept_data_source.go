package kable

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/redradrat/kable/pkg/concepts"
	"github.com/redradrat/kable/pkg/repositories"
	"reflect"
)

var ConceptDataSource = func() *schema.Resource {
	out := schema.Resource{
		Schema: map[string]*schema.Schema{
			"repo": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The repository to get the resource from",
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The concept identifier.",
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
		Read:        ConceptRead,
		Description: "The concept data source allows for rendering a concept from a repository",
	}

	return &out
}

var ConceptRead = func(d *schema.ResourceData, m interface{}) error {
	conceptPath := d.Get("id").(string)
	targetType := d.Get("target_type").(string)
	rmap := d.Get("repo").(map[string]interface{})
	r := repositories.Repository{
		Name: rmap["name"].(string),
		GitRepository: repositories.GitRepository{
			URL: rmap["url"].(string),
		},
	}
	user := rmap["username"].(string)
	if user != "" {
		r.Username = strPtr(user)
	}
	ref := rmap["ref"].(string)
	if ref != "" {
		r.GitRef = ref
	}

	// Update the registry
	if err := repositories.UpdateRegistry(repositories.AddRepository(r)); err != nil {
		return err
	}

	conceptIdentifier := concepts.NewConceptIdentifier(conceptPath, r.Name)
	avs, err := assertValues(d)
	if err != nil {
		return err
	}

	if err := renderConcept(d, conceptIdentifier.String(), avs, targetType, false); err != nil {
		return err
	}

	return nil
}

func renderConcept(d *schema.ResourceData, id string, avs *concepts.RenderValues, target string, local bool) error {
	rdr, err := concepts.RenderConcept(id, avs, concepts.TargetType(target), concepts.RenderOpts{
		Local:           local,
		WriteRenderInfo: false,
		Single:          true,
	})
	if err != nil {
		return err
	}

	if err := d.Set("rendered", rdr.PrintFiles()); err != nil {
		return err
	}

	return nil
}

func assertValues(d *schema.ResourceData) (*concepts.RenderValues, error) {
	avs := concepts.RenderValues{}
	values := d.Get("values").(map[string]interface{})
	for k, v := range values {
		switch assertedType := v.(type) {
		case bool:
			avs[k] = concepts.BoolValueType(assertedType)
		case string:
			avs[k] = concepts.StringValueType(assertedType)
		case int:
			avs[k] = concepts.IntValueType(assertedType)
		case map[string]interface{}:
			avs[k] = concepts.MapValueType(assertedType)
		default:
			return nil, fmt.Errorf("unsupported type in values: %s", reflect.TypeOf(v).String())
		}
	}
	return &avs, nil
}

func strPtr(s string) *string {
	return &s
}
