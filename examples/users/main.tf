terraform {
  required_providers {
    thebastion = {
      source = "hashicorp.com/edu/thebastion"
    }
  }
}

provider "thebastion" {}

data "thebastion_users" "all" {}

output "edu_users" {
  value = data.thebastion_users.all
}
