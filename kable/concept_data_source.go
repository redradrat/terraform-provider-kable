package kable

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/redradrat/kable/pkg/concepts"
	"github.com/redradrat/kable/pkg/repositories"
)

var ConcepRepoDataSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"url": {
		Type:     schema.TypeString,
		Required: true,
	},
	"ref": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"username": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"password": {
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	},
}

var ConceptDataSource = func() *schema.Resource {
	out := schema.Resource{
		Schema: map[string]*schema.Schema{
			"repo": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: ConcepRepoDataSchema,
				},
				Optional:    true,
				Description: "The repository to get the resource from",
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The concept identifier.",
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

	out := rdr.PrintFiles()
	if err := d.Set("rendered", out); err != nil {
		return err
	}

	d.SetId(hash(out))

	return nil
}

func assertValues(d *schema.ResourceData) (*concepts.RenderValues, error) {
	avs := concepts.RenderValues{}
	values := d.Get("inputs").(*schema.Set).List()
	for _, v := range values {
		vmap := v.(map[string]interface{})
		str := vmap["value"].(string)
		var val interface{}
		if json.Valid([]byte(str)) {
			if err := json.Unmarshal([]byte(str), &val); err != nil {
				return nil, err
			}
		} else {
			val = str
		}
		name := vmap["name"].(string)
		switch assertedType := val.(type) {
		case bool:
			avs[name] = concepts.BoolValueType(assertedType)
		case string:
			avs[name] = concepts.StringValueType(assertedType)
		case int:
			avs[name] = concepts.IntValueType(assertedType)
		case map[string]interface{}:
			avs[name] = concepts.MapValueType(assertedType)
		default:
			return nil, fmt.Errorf("unsupported type in values: %s", reflect.TypeOf(v).String())
		}
	}
	return &avs, nil
}

func strPtr(s string) *string {
	return &s
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
