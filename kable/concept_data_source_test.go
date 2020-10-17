package kable

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const renderedContent = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  minReadySeconds: 10
  replicas: 2
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      name: test
  template:
    metadata:
      labels:
        name: test
    spec:
      containers:
      - image: grafana/grafana
        imagePullPolicy: IfNotPresent
        name: grafana
        ports:
        - containerPort: 10330
          name: ui
---
apiVersion: v1
kind: Service
metadata:
  labels:
    name: test
  name: test
spec:
  ports:
  - name: grafana-ui
    port: 10330
    targetPort: 10330
  selector:
    name: test
`

func TestConcept(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{

			{
				Config: testConceptConfig(),
				Check: func(s *terraform.State) error {
					render := s.RootModule().Outputs["rendered"].Value.(string)
					if render != renderedContent {
						return fmt.Errorf("unexpected render content")
					}
					return nil
				},
			},
		},
	})
}

func testConceptConfig() string {
	return `
data "kable_concept" "bla" {
  repo {
    name = "demo"
    url = "https://github.com/redradrat/demo-concepts"
  }
  id = "apps/grafana@demo"
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
	value = data.kable_concept.bla.rendered
}
`
}
