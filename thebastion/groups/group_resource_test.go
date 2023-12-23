package groups_test

import (
	"terraform-provider-thebastion/thebastion"
	"terraform-provider-thebastion/thebastion/groups"
	"terraform-provider-thebastion/thebastion/tests"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTheBastionGroup_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	owner := "poweruser"
	algo := "rsa"
	size := 2048

	server1_host := "127.0.0.1"
	server1_port := int64(22)
	server1_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server1_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server1 := groups.ServerModel{
		Host:        types.StringValue(server1_host),
		Port:        types.Int64Value(server1_port),
		User:        types.StringValue(server1_user),
		UserComment: types.StringValue(server1_user_comment),
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
				Config:  tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
				Check:   tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
				Destroy: false,
			},
		},
	})
}

func TestAccTheBastionGroup_multiple_servers(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	owner := "poweruser"
	algo := "rsa"
	size := 2048

	server1_host := "127.0.0.1"
	server1_port := int64(22)
	server1_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server1_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server2_host := "127.0.0.2"
	server2_port := int64(22)
	server2_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server2_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server1 := groups.ServerModel{
		Host:        types.StringValue(server1_host),
		Port:        types.Int64Value(server1_port),
		User:        types.StringValue(server1_user),
		UserComment: types.StringValue(server1_user_comment),
	}

	server2 := groups.ServerModel{
		Host:        types.StringValue(server2_host),
		Port:        types.Int64Value(server2_port),
		User:        types.StringValue(server2_user),
		UserComment: types.StringValue(server2_user_comment),
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
				Config:  tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.ServerModel{server1, server2}),
				Check:   tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.ServerModel{server1, server2}),
				Destroy: false,
			},
		},
	})
}

func TestAccTheBastionGroup_update_servers(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	owner := "poweruser"
	algo := "rsa"
	size := 2048

	server1_host := "127.0.0.1"
	server1_port := int64(22)
	server1_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server1_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server2_host := "127.0.0.2"
	server2_port := int64(22)
	server2_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server2_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server1 := groups.ServerModel{
		Host:        types.StringValue(server1_host),
		Port:        types.Int64Value(server1_port),
		User:        types.StringValue(server1_user),
		UserComment: types.StringValue(server1_user_comment),
	}

	server2 := groups.ServerModel{
		Host:        types.StringValue(server2_host),
		Port:        types.Int64Value(server2_port),
		User:        types.StringValue(server2_user),
		UserComment: types.StringValue(server2_user_comment),
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
				Config: tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
				Check:  tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
			},
			{
				Config: tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.ServerModel{server2}),
				Check:  tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.ServerModel{server2}),
			},
		},
	})
}

func TestAccTheBastionGroup_update_name(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	owner := "poweruser"
	algo := "rsa"
	size := 2048

	server1_host := "127.0.0.1"
	server1_port := int64(22)
	server1_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server1_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server1 := groups.ServerModel{
		Host:        types.StringValue(server1_host),
		Port:        types.Int64Value(server1_port),
		User:        types.StringValue(server1_user),
		UserComment: types.StringValue(server1_user_comment),
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
				Config: tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
				Check:  tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
			},
			{
				Config: tests.TestAccTheBastionGroupResource(resourceName, name+"-updated", owner, algo, size, []groups.ServerModel{server1}),
				Check:  tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name+"-updated", owner, algo, size, []groups.ServerModel{server1}),
			},
		},
	})
}

func TestAccTheBastionGroup_update_owner(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	owner := "poweruser"
	owner_updated := "healthcheck"
	algo := "rsa"
	size := 2048

	server1_host := "127.0.0.1"
	server1_port := int64(22)
	server1_user := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	server1_user_comment := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	server1 := groups.ServerModel{
		Host:        types.StringValue(server1_host),
		Port:        types.Int64Value(server1_port),
		User:        types.StringValue(server1_user),
		UserComment: types.StringValue(server1_user_comment),
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
				Config: tests.TestAccTheBastionGroupResource(resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
				Check:  tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner, algo, size, []groups.ServerModel{server1}),
			},
			{
				Config: tests.TestAccTheBastionGroupResource(resourceName, name, owner_updated, algo, size, []groups.ServerModel{server1}),
				Check:  tests.TestAccCheckTheBastionGroupValues("thebastion_group."+resourceName, name, owner_updated, algo, size, []groups.ServerModel{server1}),
			},
		},
	})
}
