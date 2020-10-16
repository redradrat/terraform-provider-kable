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
    path = "${path.module}/../kable/test/concept"
    inputs {
        name = "instanceName"
        value = "test"
    }
    inputs {
        name = "nameSelection"
        value = "Option 1"
    }
}


resource "local_file" "foo" {
    content     = data.kable_local_concept.test.rendered
    filename = "${path.module}/foo.yaml"
}
