package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"terraform-provider-thebastion/thebastion/clients"
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
	checks := []resource.TestCheckFunc{}

	// Check the list of ingress_keys
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "ingress_keys.#", fmt.Sprint(len(ingress_keys))))
	for i, ingress_key := range ingress_keys {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "ingress_keys."+fmt.Sprint(i), ingress_key))
	}

	// Check the other attributes
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "name", name))
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "uid", uid))
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "is_active", is_active))

	return resource.ComposeAggregateTestCheckFunc(checks...)
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

func TestAccCheckTheBastionGroupValues(resourceName, name string, owners []string, algo string, size int, list_server []clients.ServerModel) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}

	// Check the list of owners
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "owners.#", fmt.Sprint(len(owners))))
	for i, owner := range owners {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, "owners."+fmt.Sprint(i), owner))
	}

	// Check the list of servers
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "servers.#", fmt.Sprint(len(list_server))))
	for i, server := range list_server {
		prefixKey := "servers." + fmt.Sprint(i) + "."
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, prefixKey+"host", server.Host.ValueString()),
			resource.TestCheckResourceAttr(resourceName, prefixKey+"user", server.User.ValueString()),
			resource.TestCheckResourceAttr(resourceName, prefixKey+"port", fmt.Sprint(server.Port)),
			resource.TestCheckResourceAttr(resourceName, prefixKey+"user_comment", server.UserComment.ValueString()),
		)
	}

	// Check the other attributes
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "name", name))
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "algo", algo))
	checks = append(checks, resource.TestCheckResourceAttr(resourceName, "size", fmt.Sprint(size)))

	return resource.ComposeAggregateTestCheckFunc(checks...)
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

func TestAccTheBastionGroupResource(resourceName string, name string, owners []string, algo string, size int, list_server []clients.ServerModel) string {
	var servers string
	for _, server := range list_server {
		servers += fmt.Sprintf(`
		{
			host = %s
			user = %s
			port = %d
			user_comment = %s
		},`, server.Host.String(), server.User.String(), server.Port.ValueInt64(), server.UserComment.String())
	}

	return fmt.Sprintf(`
	resource "thebastion_group" "%s" {
		name = "%s"
		owners = ["%s"]
		algo = "%s"
		size = %d
		servers = [%s]
	}`, resourceName, name, strings.Join(owners, "\",\""), algo, size, servers)
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

	_, err := GetClient()
	require.Nil(err, "Cannot connect to TheBastion.")
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
