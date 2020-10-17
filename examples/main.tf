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


data "kable_concept" "demo" {
  repo {
    name = "demo"
    url = "https://github.com/demo/repository"
    username = "$env.REPO_USER"
    password = "$env.REPO_PASS"
  }
  id = "concept@demo"
  inputs {
    name = "name"
    value = "test"
  }
  inputs {
    name = "namespace"
    value = "default"
  }
}

