package kable

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

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
			}, "sensitive_inputs": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The key/value list of sensitive inputs to use.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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
		ReadContext: ConceptRead,
		Description: "The concept data source allows for rendering a concept from a repository",
	}

	return &out
}
var ConceptRead = func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conceptPath := d.Get("id").(string)
	targetType := d.Get("target_type").(string)

	var mods []repositories.RegistryModification

	for _, repo := range d.Get("repo").(*schema.Set).List() {
		rmap := repo.(map[string]interface{})
		name := rmap["name"].(string)
		url := rmap["url"].(string)
		ref := rmap["ref"].(string)
		r := repositories.Repository{
			Name: name,
			GitRepository: repositories.GitRepository{
				URL: url,
			},
		}
		if ref != "" {
			r.GitRef = ref
		}
		addRepoMod, err := repositories.AddRepository(r)
		if err != nil {
			return diag.FromErr(err)
		}
		mods = append(mods, addRepoMod)

		// Let's see if we have some auth defined
		user := rmap["username"].(string)
		pass := rmap["password"].(string)
		if user != "" && pass != "" {
			authmod, err := repositories.StoreRepoAuth(url, repositories.AuthPair{Username: user, Password: pass})
			if err != nil {
				return diag.FromErr(err)
			}
			mods = append(mods, authmod)
		}
	}

	// Update the registry
	if err := repositories.UpdateRegistry(mods...); err != nil {
		return diag.FromErr(err)
	}

	conceptIdentifier := concepts.ConceptIdentifier(conceptPath)
	avs, err := assertValues(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Now let's render our concept
	if err := renderConcept(d, conceptIdentifier.String(), avs, targetType, false); err != nil {
		return diag.FromErr(err)
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
	sensitiveValues := d.Get("sensitive_inputs").(*schema.Set).List()
	values = append(values, sensitiveValues...)
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
