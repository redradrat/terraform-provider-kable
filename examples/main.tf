terraform {
  required_providers {
    kable = {
      source = "redradrat/kable"
    }
  }
}

provider "kable" {
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
    sensitive_inputs {
        name = "name"
        value = "test"
    }
}

