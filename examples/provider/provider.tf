# Configuration authentication required for TheBastion
provider "thebastion" {
  host = "host"
  username = "usernameBastion"
  path_private_key = "/Users/username/.ssh/id_ed25519"
  path_known_host = "/Users/username/.ssh/known_hosts"
}
