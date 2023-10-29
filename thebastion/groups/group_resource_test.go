package groups_test

import (
	"terraform-provider-thebastion/thebastion"
	"terraform-provider-thebastion/thebastion/groups"
	"terraform-provider-thebastion/thebastion/tests"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTheBastionGroup_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	owner := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	algo := "rsa"
	size := 2048

	server1_host := "127.0.0.1"
	server1_port := 22
	server1_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server1_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server1 := groups.Server{
		Host:    server1_host,
		Port:    int64(server1_port),
		User:    server1_user,
		Comment: server1_comment,
	}

	thebastionServer := thebastion.New()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"thebastion": providerserver.NewProtocol6WithError(thebastionServer),
		},
		CheckDestroy: tests.TestAccCheckTheBastionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:  tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.Server{server1}),
				Check:   tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.Server{server1}),
				Destroy: true,
			},
		},
	})
}
