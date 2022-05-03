terraform {

  required_providers {
    ranger = {
      source  = "vranyes.com/apache/ranger"
      version = "~> 0.0.0"
    }
  }
}

provider "ranger" {
  host            = "<HOST>"
  username        = "admin"
  password        = "<pw>"
  skip_ssl_verify = true
}

data "ranger_policy" "test" {
  id = "22"
}

resource "ranger_policy" "test" {
  name         = "josiahs test policy"
  labels       = ["/terraform"]
  service_type = "hbase"
  service      = "vranyj2hdp3_hbase"
  resource {
    key    = "column-family"
    values = ["*"]
  }
  resource {
    key      = "column"
    values   = ["*"]
    excludes = true
  }
  resource {
    key       = "table"
    values    = ["*"]
    recursive = true
  }
  policy {
    accesses = ["read", "write", "create", "admin"]
    users    = ["vranyj2"]
    groups   = ["admins"]
  }
}

output "name" {
  value = data.ranger_policy.test.name
}