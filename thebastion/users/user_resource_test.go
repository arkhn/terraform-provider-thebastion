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

func TestAccTheBastionUser_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	// generate random uid for test
	// We exclued 9998, 9999 values setted for healthcheck and poweruser
	// users at start of thebastion
	rUid := int64(acctest.RandIntRange(2000, 9997))

	ingress_keys_base := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcjliyS0gOlGrxz0bX0S6GV1roGW2beEiIB+/yzygXzL7vzRU3u6Ty/wODC+kABNebtgJ7TCFj387drS3A14bojFlbSlS+r9bdToczfc0ZxwV89ToEGkw4hWIsTSw2ADg9aTIDclAZjNtE+SQUZLSS1gKJSHKah4SWaMf7CSHy7zKg4Q70qHEXJ+UCPfR30glX7joH5kny81aY9vRtRQKs6/RbG8Zd2CoxBkNAYA2k9NPVKEv3eUhiwkK+c1Zf9L5Fk2mW1jhvOwQ4auvZdV/mh/mY5uWqV2Q7KjhpucnVVgv87Uv6drL2lvQyDOvl1G03ab+rXS7eKD3aX1MkphxCrSsNaG4lTT0NB72Wa64CrCHGMcqPrdAhHkRnze/XdmXW7FOlo+nmLPRBZlBME+XT9yyQFNxksJpTAZEK33Xwccoq9PwqPsOFIHPS8PiVifQMarLXonlCz++wzoFEsdYCxdvU/jJmjBvsBcFXV+V5whtOc9JGAJ6JrtnEJJd774c="}
	ingress_keys_update := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDDCkSqTiaV8+KU6s9DEioMK5C99pVPpdf5Sx3VQPr6eQFqQj+luP8SLeUMuQI0Q+S/2mY0QvDffF0pcfTVS6VxJ160aA525kLrkFnKrORHTnAifOObjrvSuUHTjrrS41RgjdHFgP0fhdA5I+CRuAkm4KwYAdQt/CifhaF2W0TFSaebe/jmRgbTuCVblJBdPlW8A9rt7VuF7tDvYclXhYScXKqAXu+kGUl9Ts/1NK9G1XuzGyoufPOQtrprzrrIx7+J0gkYPkgz5K7yvNj2HKfiCMoB+ifb0PrWAis/ois+7nTG5FYFTdCvs7Nt+HOlJ174dmQJWV6PpExZSsO2SdfgNJM+PGU3rcvCl1Z/cRrGDnS8OKNlu6qhrNgLZLvTedS45UO0R59xqaHoJ8Ge6LyDRW4YIpNEmLl9Mn+38njSdkU+2rWzqBCt9uVBXVfTEL9t8ceHJzPkUBghNWaKjAYdjE27h2cy6gHkDLVIbjq0nthS2RMKa2Txho4qpVXjfSE="}
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
				Config: tests.TestAccTheBastionUserResource(resourceName, rUid, name, ingress_keys_update),
				Check:  tests.TestAccCheckTheBastionUserValues("thebastion_user."+resourceName, fmt.Sprint(rUid), name, "1", ingress_keys_update),
			},
		},
	})
}

func TestAccTheBastionUser_update_name(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	nameUpdate := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	// generate random uid for test
	// We exclued 9998, 9999 values setted for healthcheck and poweruser
	// users at start of thebastion
	rUid := int64(acctest.RandIntRange(2000, 9997))
	ingress_keys := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDcjliyS0gOlGrxz0bX0S6GV1roGW2beEiIB+/yzygXzL7vzRU3u6Ty/wODC+kABNebtgJ7TCFj387drS3A14bojFlbSlS+r9bdToczfc0ZxwV89ToEGkw4hWIsTSw2ADg9aTIDclAZjNtE+SQUZLSS1gKJSHKah4SWaMf7CSHy7zKg4Q70qHEXJ+UCPfR30glX7joH5kny81aY9vRtRQKs6/RbG8Zd2CoxBkNAYA2k9NPVKEv3eUhiwkK+c1Zf9L5Fk2mW1jhvOwQ4auvZdV/mh/mY5uWqV2Q7KjhpucnVVgv87Uv6drL2lvQyDOvl1G03ab+rXS7eKD3aX1MkphxCrSsNaG4lTT0NB72Wa64CrCHGMcqPrdAhHkRnze/XdmXW7FOlo+nmLPRBZlBME+XT9yyQFNxksJpTAZEK33Xwccoq9PwqPsOFIHPS8PiVifQMarLXonlCz++wzoFEsdYCxdvU/jJmjBvsBcFXV+V5whtOc9JGAJ6JrtnEJJd774c="}

	thebastionServer := thebastion.New()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { tests.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"thebastion": providerserver.NewProtocol6WithError(thebastionServer),
		},
		CheckDestroy: tests.TestAccCheckTheBastionUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: tests.TestAccTheBastionUserResource(resourceName, rUid, name, ingress_keys),
				Check:  tests.TestAccCheckTheBastionUserValues("thebastion_user."+resourceName, fmt.Sprint(rUid), name, "1", ingress_keys),
			},
			{
				Config: tests.TestAccTheBastionUserResource(resourceName, rUid, nameUpdate, ingress_keys),
				Check:  tests.TestAccCheckTheBastionUserValues("thebastion_user."+resourceName, fmt.Sprint(rUid), nameUpdate, "1", ingress_keys),
			},
		},
	})
}
