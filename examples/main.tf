terraform {
  required_providers {
    kable = {
      source = "github.com/redradrat/kable"
    }
  }
}

provider "kable" {
    version = "0.1"
}

data "kable_local_concept" "test" {
    path = "/Users/ralphkuehnert/Development/hilti/concepts/apps/observability/glowroot"
    target_type = "yaml"
    values = {
        instanceName: "test"
        name: "blabla"
        namespace: "sdalfjlasd"
        # ingressAnnotations: {} <- map not yet possible
        domain: "sadfdsaf"
        ingressCertSecret: "dsifoasjdi"
    }
}