terraform {
  required_providers {
    thebastion = {
      source = "hashicorp.com/ovh/thebastion"
    }
  }
}

provider "thebastion" {}

data "thebastion_users" "all" {}