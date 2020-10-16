package kable

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const localRenderedContent = `apiVersion: v1
kind: Test
metadata:
  name: test
  namespace: Option 1
`

var testProviders = map[string]*schema.Provider{
	"kable": Provider(),
}

func TestLocalConcept(t *testing.T) {
	abs, err := filepath.Abs("./test/concept")
	if err != nil {
		t.Error(err)
	}
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{

			{
				Config: testLocalConceptConfig(abs),
				Check: func(s *terraform.State) error {
					render := s.RootModule().Outputs["rendered"].Value.(string)
					if render != localRenderedContent {
						return fmt.Errorf("unexpected render content")
					}
					return nil
				},
			},
		},
	})
}

func testLocalConceptConfig(path string) string {
	return fmt.Sprintf(`
data "kable_local_concept" "test" {
	path = "%s"
	inputs {
		name = "instanceName"
		value = "test"
	}
	inputs {
		name = "nameSelection"
		value = "Option 1"
	}
}

output "rendered" {
	value = data.kable_local_concept.test.rendered
}
`, path)
}
