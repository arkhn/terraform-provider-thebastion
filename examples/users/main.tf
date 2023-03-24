terraform {
  required_providers {
    thebastion = {
      source = "hashicorp.com/ovh/thebastion"
    }
  }
}

provider "thebastion" {}

resource "thebastion_user" "test" {
    name = "test"
    uid =  6666
    ingress_keys = [
      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcjliyS0gOlGrxz0bX0S6GV1roGW2beEiIB+/yzygXzL7vzRU3u6Ty/wODC+kABNebtgJ7TCFj387drS3A14bojFlbSlS+r9bdToczfc0ZxwV89ToEGkw4hWIsTSw2ADg9aTIDclAZjNtE+SQUZLSS1gKJSHKah4SWaMf7CSHy7zKg4Q70qHEXJ+UCPfR30glX7joH5kny81aY9vRtRQKs6/RbG8Zd2CoxBkNAYA2k9NPVKEv3eUhiwkK+c1Zf9L5Fk2mW1jhvOwQ4auvZdV/mh/mY5uWqV2Q7KjhpucnVVgv87Uv6drL2lvQyDOvl1G03ab+rXS7eKD3aX1MkphxCrSsNaG4lTT0NB72Wa64CrCHGMcqPrdAhHkRnze/XdmXW7FOlo+nmLPRBZlBME+XT9yyQFNxksJpTAZEK33Xwccoq9PwqPsOFIHPS8PiVifQMarLXonlCz++wzoFEsdYCxdvU/jJmjBvsBcFXV+V5whtOc9JGAJ6JrtnEJJd774c="
    ]
}

data "thebastion_users" "all" {}
