resource "thebastion_user" "example" {
    name = "example"
    uid = 88888
    ingress_keys = [
        "... ingress key 1 ...",
        "... ingress key 2 ..."
    ]
}