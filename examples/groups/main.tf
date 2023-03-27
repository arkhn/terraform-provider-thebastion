terraform {
  required_providers {
    thebastion = {
      source = "hashicorp/arkhn/thebastion"
    }
  }
}

provider "thebastion" {}

data "thebastion_groups" "all" {}
