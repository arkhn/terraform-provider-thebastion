---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "thebastion_user Resource - thebastion"
subcategory: ""
description: |-
  Manage an user.
---

# thebastion_user (Resource)

Manage an user.

## Example Usage

```terraform
resource "thebastion_user" "example" {
    name = "example"
    uid = 88888
    ingress_keys = [
        "... ingress key 1 ...",
        "... ingress key 2 ..."
    ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ingress_keys` (List of String) List of ingress keys of users.
- `name` (String) Name of user. Used as an unique identifier by TheBastion.
- `uid` (Number) UID of user.

### Read-Only

- `id` (String) ID of resource. Required by terraform-plugin-testing
- `is_active` (Number) Is the user active.

## Import

Import is supported using the following syntax:

```shell
# Order can be imported by specifying the numeric identifier.
terraform import thebastion_user.example example
```