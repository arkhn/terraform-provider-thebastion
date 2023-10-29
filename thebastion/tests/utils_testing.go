package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"terraform-provider-thebastion/thebastion/clients"
	"terraform-provider-thebastion/thebastion/groups"
	"terraform-provider-thebastion/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

var lock = &sync.Mutex{}

var clientInstance *clients.Client

func GetClient() (*clients.Client, error) {
	if clientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if clientInstance == nil {
			// Create the new instance of client
			host := os.Getenv("THEBASTION_HOST")
			username := os.Getenv("THEBASTION_USERNAME")
			path_known_host := os.Getenv("THEBASTION_PATH_KNOWN_HOST")
			path_private_key := os.Getenv("THEBASTION_PATH_PRIVATE_KEY")

			var err error
			clientInstance, err = clients.NewClient(host, username, path_private_key, path_known_host)
			if err != nil {
				return clientInstance, fmt.Errorf("cannot connect to TheBastion: %s", err.Error())
			}
		}
	}

	return clientInstance, nil
}

func TestAccCheckTheBastionUserValues(resourceName, uid, name, is_active string, ingress_keys []string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "is_active", is_active),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "uid", uid),
		resource.TestCheckResourceAttr(resourceName, "ingress_keys.#", fmt.Sprint(len(ingress_keys))),
		resource.TestCheckResourceAttr(resourceName, "ingress_keys.0", ingress_keys[0]),
	)
}

func TestAccCheckTheBastionUsersValues(resourceName string, i int, uid, name, is_active string, ingress_keys []string) resource.TestCheckFunc {
	prefixKey := "users." + fmt.Sprint(i) + "."
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, prefixKey+"is_active", is_active),
		resource.TestCheckResourceAttr(resourceName, prefixKey+"name", name),
		resource.TestCheckResourceAttr(resourceName, prefixKey+"uid", uid),
		resource.TestCheckResourceAttr(resourceName, prefixKey+"ingress_keys.#", fmt.Sprint(len(ingress_keys))),
		resource.TestCheckResourceAttr(resourceName, prefixKey+"ingress_keys.0", ingress_keys[0]),
	)
}

func TestAccCheckTheBastionGroupValues(resourceName, name, owner, algo string, size int, list_server []groups.Server) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "owner", owner),
		resource.TestCheckResourceAttr(resourceName, "algo", algo),
		resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprint(size)),
		resource.TestCheckResourceAttr(resourceName, "servers.#", fmt.Sprint(len(list_server))),
		resource.TestCheckResourceAttr(resourceName, "servers.0.host", list_server[0].Host),
		resource.TestCheckResourceAttr(resourceName, "servers.0.user", list_server[0].User),
		resource.TestCheckResourceAttr(resourceName, "servers.0.port", fmt.Sprint(list_server[0].Port)),
		resource.TestCheckResourceAttr(resourceName, "servers.0.comment", list_server[0].Comment),
	)
}

// TestAccTheBastionUserResource returns an configuration for an user with the provided configuration
func TestAccTheBastionUserResource(resourceName string, uid int64, name string, ingress_keys []string) string {
	return fmt.Sprintf(`
	resource "thebastion_user" "%s" {
		uid = %8d
		name = "%s"
		ingress_keys = ["%s"]
	}`, resourceName, uid, name, strings.Join(ingress_keys, "\",\""))
}

func TestAccTheBastionGroupResource(resourceName string, name string, owner string, algo string, size int, list_server []groups.Server) string {
	var servers string
	for _, server := range list_server {
		servers += fmt.Sprintf(`
		{
			host = "%s"
			user = "%s"
			port = %d
			comment = "%s"
		},`, server.Host, server.User, server.Port, server.Comment)
	}

	return fmt.Sprintf(`
	resource "thebastion_group" "%s" {
		name = "%s"
		owner = "%s"
		algo = "%s"
		size = %d
		servers = [%s]
	}`, resourceName, name, owner, algo, size, servers)
}

func TestAccTheBastionUserDataSource(exampleResource string) string {
	return exampleResource + `
	data "thebastion_users" "all" {}
	`
}

func TestAccTheBastionGroupsDataSource(exampleResource string) string {
	return exampleResource + `
	data "thebastion_groups" "all" {}
	`
}

// TestAccPreCheck validates the necessary test API keys exist
// in the testing environment
func TestAccPreCheck(t *testing.T) {
	// Terminate test if fail
	require := require.New(t)

	host := os.Getenv("THEBASTION_HOST")
	username := os.Getenv("THEBASTION_USERNAME")
	path_known_host := os.Getenv("THEBASTION_PATH_KNOWN_HOST")
	path_private_key := os.Getenv("THEBASTION_PATH_PRIVATE_KEY")

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	msg_error, msg_error_detail := utils.MissingEnvMsg("host", "THEBASTION_HOST")
	require.NotEqual(host, "", msg_error, msg_error_detail)

	msg_error, msg_error_detail = utils.MissingEnvMsg("username", "THEBASTION_USERNAME")
	require.NotEqual(username, "", msg_error, msg_error_detail)

	msg_error, msg_error_detail = utils.MissingEnvMsg("path_known_host", "THEBASTION_PATH_KNOWN_HOST")
	require.NotEqual(path_known_host, "", msg_error, msg_error_detail)

	msg_error, msg_error_detail = utils.MissingEnvMsg("path_private_key", "THEBASTION_PATH_PRIVATE_KEY")
	require.NotEqual(path_private_key, "", msg_error, msg_error_detail)

	client, err := GetClient()
	require.Nil(err, "Cannot connect to TheBastion.")

	// Make sure only expected users are on TheBastion
	responseBastionAccountList, err := client.GetListAccount(context.Background())
	require.Nil(err, "Cannot get the list of account from TheBastion.")

	nbUsersTheBastion := 2
	require.Equal(len(responseBastionAccountList.Value), nbUsersTheBastion, "Unexpected users on TheBastion for testing. Please delete all users on TheBastion except poweruser and healthcheck: "+fmt.Sprint(responseBastionAccountList.Value))

	// Make sure no groups are on TheBastion
	responseBastionGroupList, err := client.GetListGroup(context.Background())
	require.Nil(err, "Cannot get the list of groups from TheBastion.")

	nbGroupsTheBastion := 0
	require.Equal(len(responseBastionGroupList.Value), nbGroupsTheBastion, "Unexpected groups on TheBastion for testing. Please delete all groups on TheBastion: "+fmt.Sprint(responseBastionGroupList.Value))
}

func TestAccCheckTheBastionUserDestroy(s *terraform.State) error {
	client, err := GetClient()

	if err != nil {
		return fmt.Errorf("cannot connect to TheBastion: %s", err.Error())

	}

	responseBastion, err := client.GetListAccount(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get the list of account from TheBastion: %s", err.Error())
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "thebastion_user" {
			continue
		}

		value, ok := responseBastion.Value[rs.Primary.Attributes["name"]]
		if ok {
			if value.Name == rs.Primary.Attributes["name"] {
				return fmt.Errorf("user (%s) still exists", rs.Primary.Attributes["name"])
			}
		}
	}

	return nil
}

func TestAccCheckTheBastionGroupDestroy(s *terraform.State) error {
	client, err := GetClient()

	if err != nil {
		return fmt.Errorf("cannot connect to TheBastion: %s", err.Error())

	}

	responseBastion, err := client.GetListGroup(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get the list of group from TheBastion: %s", err.Error())
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "thebastion_group" {
			continue
		}

		_, ok := responseBastion.Value[rs.Primary.Attributes["name"]]
		if ok {
			return fmt.Errorf("group (%s) still exists", rs.Primary.Attributes["name"])
		}
	}

	return nil
}
