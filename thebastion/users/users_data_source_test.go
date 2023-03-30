package users_test

import (
	"fmt"
	"terraform-provider-thebastion/thebastion"
	"terraform-provider-thebastion/thebastion/tests"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccTheBastionUsers_single test to get users when there is single user in TheBastion
func TestAccTheBastionUsers_single(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	// generate random uid for test
	// We exclued 9998, 9999 values setted for healthcheck and poweruser
	// users at start of thebastion
	rUid := int64(acctest.RandIntRange(2000, 9997))

	ingress_keys_base := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcjliyS0gOlGrxz0bX0S6GV1roGW2beEiIB+/yzygXzL7vzRU3u6Ty/wODC+kABNebtgJ7TCFj387drS3A14bojFlbSlS+r9bdToczfc0ZxwV89ToEGkw4hWIsTSw2ADg9aTIDclAZjNtE+SQUZLSS1gKJSHKah4SWaMf7CSHy7zKg4Q70qHEXJ+UCPfR30glX7joH5kny81aY9vRtRQKs6/RbG8Zd2CoxBkNAYA2k9NPVKEv3eUhiwkK+c1Zf9L5Fk2mW1jhvOwQ4auvZdV/mh/mY5uWqV2Q7KjhpucnVVgv87Uv6drL2lvQyDOvl1G03ab+rXS7eKD3aX1MkphxCrSsNaG4lTT0NB72Wa64CrCHGMcqPrdAhHkRnze/XdmXW7FOlo+nmLPRBZlBME+XT9yyQFNxksJpTAZEK33Xwccoq9PwqPsOFIHPS8PiVifQMarLXonlCz++wzoFEsdYCxdvU/jJmjBvsBcFXV+V5whtOc9JGAJ6JrtnEJJd774c="}
	thebastionServer := thebastion.New()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"thebastion": providerserver.NewProtocol6WithError(thebastionServer),
		},
		CheckDestroy: tests.TestAccCheckTheBastionUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: tests.TestAccTheBastionUserResource(resourceName, rUid, name, ingress_keys_base),
				Check:  tests.TestAccCheckTheBastionUserValues("thebastion_user."+resourceName, fmt.Sprint(rUid), name, "1", ingress_keys_base),
			},
			{
				Config: tests.TestAccTheBastionUserDataSource(tests.TestAccTheBastionUserResource(resourceName, rUid, name, ingress_keys_base)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.thebastion_users.all", "users.#", "3"),
					resource.TestCheckResourceAttr("data.thebastion_users.all", "users.0.uid", fmt.Sprint(rUid)),

					tests.TestAccCheckTheBastionUsersValues("data.thebastion_users.all", 0, fmt.Sprint(rUid), name, "1", ingress_keys_base),
				),
			},
		},
	})
}

// TestAccTheBastionUsers_multiple test to get users when there is multiple users in TheBastion
func TestAccTheBastionUsers_multiple(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rUid := int64(acctest.RandIntRange(2000, 9996))
	ingress_keys := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcjliyS0gOlGrxz0bX0S6GV1roGW2beEiIB+/yzygXzL7vzRU3u6Ty/wODC+kABNebtgJ7TCFj387drS3A14bojFlbSlS+r9bdToczfc0ZxwV89ToEGkw4hWIsTSw2ADg9aTIDclAZjNtE+SQUZLSS1gKJSHKah4SWaMf7CSHy7zKg4Q70qHEXJ+UCPfR30glX7joH5kny81aY9vRtRQKs6/RbG8Zd2CoxBkNAYA2k9NPVKEv3eUhiwkK+c1Zf9L5Fk2mW1jhvOwQ4auvZdV/mh/mY5uWqV2Q7KjhpucnVVgv87Uv6drL2lvQyDOvl1G03ab+rXS7eKD3aX1MkphxCrSsNaG4lTT0NB72Wa64CrCHGMcqPrdAhHkRnze/XdmXW7FOlo+nmLPRBZlBME+XT9yyQFNxksJpTAZEK33Xwccoq9PwqPsOFIHPS8PiVifQMarLXonlCz++wzoFEsdYCxdvU/jJmjBvsBcFXV+V5whtOc9JGAJ6JrtnEJJd774c="}

	resourceName2 := resourceName + "_2"
	name2 := name + "_2"
	rUid2 := rUid + 1
	ingress_keys2 := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcjliyS0gOlGrxz0bX0S6GV1roGW2beEiIB+/yzygXzL7vzRU3u6Ty/wODC+kABNebtgJ7TCFj387drS3A14bojFlbSlS+r9bdToczfc0ZxwV89ToEGkw4hWIsTSw2ADg9aTIDclAZjNtE+SQUZLSS1gKJSHKah4SWaMf7CSHy7zKg4Q70qHEXJ+UCPfR30glX7joH5kny81aY9vRtRQKs6/RbG8Zd2CoxBkNAYA2k9NPVKEv3eUhiwkK+c1Zf9L5Fk2mW1jhvOwQ4auvZdV/mh/mY5uWqV2Q7KjhpucnVVgv87Uv6drL2lvQyDOvl1G03ab+rXS7eKD3aX1MkphxCrSsNaG4lTT0NB72Wa64CrCHGMcqPrdAhHkRnze/XdmXW7FOlo+nmLPRBZlBME+XT9yyQFNxksJpTAZEK33Xwccoq9PwqPsOFIHPS8PiVifQMarLXonlCz++wzoFEsdYCxdvU/jJmjBvsBcFXV+V5whtOc9JGAJ6JrtnEJJd774c="}

	thebastionServer := thebastion.New()
	userConfig := tests.TestAccTheBastionUserResource(resourceName, rUid, name, ingress_keys) + tests.TestAccTheBastionUserResource(resourceName2, rUid2, name2, ingress_keys2)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"thebastion": providerserver.NewProtocol6WithError(thebastionServer),
		},
		CheckDestroy: tests.TestAccCheckTheBastionUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: userConfig,
				Check:  tests.TestAccCheckTheBastionUserValues("thebastion_user."+resourceName, fmt.Sprint(rUid), name, "1", ingress_keys),
			},
			{
				Config: tests.TestAccTheBastionUserDataSource(userConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.thebastion_users.all", "users.#", "4"),
					resource.TestCheckResourceAttr("data.thebastion_users.all", "users.0.uid", fmt.Sprint(rUid)),
					resource.TestCheckResourceAttr("data.thebastion_users.all", "users.1.uid", fmt.Sprint(rUid2)),

					tests.TestAccCheckTheBastionUsersValues("data.thebastion_users.all", 0, fmt.Sprint(rUid), name, "1", ingress_keys),
					tests.TestAccCheckTheBastionUsersValues("data.thebastion_users.all", 1, fmt.Sprint(rUid2), name2, "1", ingress_keys2),
				),
			},
		},
	})
}

// TestAccTheBastionUsers_empty test to get users when there is no users in TheBastion
func TestAccTheBastionUsers_empty(t *testing.T) {
	thebastionServer := thebastion.New()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"thebastion": providerserver.NewProtocol6WithError(thebastionServer),
		},
		Steps: []resource.TestStep{
			{
				Config: tests.TestAccTheBastionUserDataSource(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.thebastion_users.all", "users.#", "2"),
				),
			},
		},
	})
}
