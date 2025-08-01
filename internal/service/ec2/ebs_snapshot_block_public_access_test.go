// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfknownvalue "github.com/hashicorp/terraform-provider-aws/internal/acctest/knownvalue"
	tfstatecheck "github.com/hashicorp/terraform-provider-aws/internal/acctest/statecheck"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccEC2EBSSnapshotBlockPublicAccess_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		acctest.CtBasic: testAccEC2EBSSnapshotBlockPublicAccess_basic,
		"Identity":      testAccEC2EBSEBSSnapshotBlockPublicAccess_IdentitySerial,
	}

	acctest.RunSerialTests1Level(t, testCases, 0)
}

func testAccEC2EBSSnapshotBlockPublicAccess_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ebs_snapshot_block_public_access.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		WorkingDir:               "/tmp",
		CheckDestroy:             testAccCheckEBSSnapshotBlockPublicAccessDestroy(ctx),
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccEBSSnapshotBlockPublicAccess_basic(string(types.SnapshotBlockPublicAccessStateBlockAllSharing)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, names.AttrState, "block-all-sharing"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName: resourceName,
				Config:       testAccEBSSnapshotBlockPublicAccess_basic(string(types.SnapshotBlockPublicAccessStateBlockNewSharing)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, names.AttrState, "block-new-sharing"),
				),
			},
		},
	})
}

func testAccEC2EBSEBSSnapshotBlockPublicAccess_Identity_ExistingResource(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ebs_snapshot_block_public_access.test"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		PreCheck:     func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:   acctest.ErrorCheck(t, names.EC2ServiceID),
		CheckDestroy: testAccCheckEBSSnapshotBlockPublicAccessDestroy(ctx),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "5.100.0",
					},
				},
				Config: testAccEBSSnapshotBlockPublicAccess_basic("block-all-sharing"),
				ConfigStateChecks: []statecheck.StateCheck{
					tfstatecheck.ExpectNoIdentity(resourceName),
				},
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "6.0.0",
					},
				},
				Config: testAccEBSSnapshotBlockPublicAccess_basic("block-all-sharing"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity(resourceName, map[string]knownvalue.Check{
						names.AttrAccountID: knownvalue.Null(),
						names.AttrRegion:    knownvalue.Null(),
					}),
				},
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
				Config:                   testAccEBSSnapshotBlockPublicAccess_basic("block-all-sharing"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity(resourceName, map[string]knownvalue.Check{
						names.AttrAccountID: tfknownvalue.AccountID(),
						names.AttrRegion:    knownvalue.StringExact(acctest.Region()),
					}),
				},
			},
		},
	})
}

func testAccCheckEBSSnapshotBlockPublicAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Client(ctx)
		input := ec2.GetSnapshotBlockPublicAccessStateInput{}
		response, err := conn.GetSnapshotBlockPublicAccessState(ctx, &input)
		if err != nil {
			return err
		}

		if response.State != types.SnapshotBlockPublicAccessStateUnblocked {
			return fmt.Errorf("EBS encryption by default is not in expected state (%s)", types.SnapshotBlockPublicAccessStateUnblocked)
		}
		return nil
	}
}

func testAccEBSSnapshotBlockPublicAccess_basic(state string) string {
	return fmt.Sprintf(`
resource "aws_ebs_snapshot_block_public_access" "test" {
  state = %[1]q
}
`, state)
}
