package main

import (
	"context"
	"terraform-provider-thebastion/thebastion"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name thebastion

func main() {
	err := providerserver.Serve(context.Background(), thebastion.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/arkhn/thebastion. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp/arkhn/thebastion",
	})
	if err != nil {
		panic(err)
	}
}
