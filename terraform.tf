terraform {

  required_providers {
    ranger = {
      source  = "vranyes.com/apache/ranger"
      version = "~> 0.0.0"
    }
  }
}

provider "ranger" {
  host = "<HOST>"
  username = "admin"
  password = "<pw>"
}

data "ranger_policy" "test" {
  id = "22"
}

output "name" {
  value = data.ranger_policy.test.name
}