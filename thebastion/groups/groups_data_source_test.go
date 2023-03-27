package groups_test

import (
	"terraform-provider-thebastion/thebastion"
	"terraform-provider-thebastion/thebastion/tests"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTheBastionGroups_empty(t *testing.T) {
	thebastionServer := thebastion.New()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"thebastion": providerserver.NewProtocol6WithError(thebastionServer),
		},
		Steps: []resource.TestStep{
			{
				Config: tests.TestAccTheBastionGroupsDataSource(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.thebastion_groups.all", "groups.#", "0"),
				),
			},
		},
	})
}
