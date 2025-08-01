// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package servicecatalog

import (
	"context"
	"unique"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	inttypes "github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/internal/vcr"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*inttypes.ServicePackageFrameworkDataSource {
	return []*inttypes.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*inttypes.ServicePackageFrameworkResource {
	return []*inttypes.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*inttypes.ServicePackageSDKDataSource {
	return []*inttypes.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceConstraint,
			TypeName: "aws_servicecatalog_constraint",
			Name:     "Constraint",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  dataSourceLaunchPaths,
			TypeName: "aws_servicecatalog_launch_paths",
			Name:     "Launch Paths",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  dataSourcePortfolio,
			TypeName: "aws_servicecatalog_portfolio",
			Name:     "Portfolio",
			Tags:     unique.Make(inttypes.ServicePackageResourceTags{}),
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  dataSourcePortfolioConstraints,
			TypeName: "aws_servicecatalog_portfolio_constraints",
			Name:     "Portfolio Constraints",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  dataSourceProduct,
			TypeName: "aws_servicecatalog_product",
			Name:     "Product",
			Tags:     unique.Make(inttypes.ServicePackageResourceTags{}),
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  dataSourceProvisioningArtifacts,
			TypeName: "aws_servicecatalog_provisioning_artifacts",
			Name:     "Provisioning Artifacts",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*inttypes.ServicePackageSDKResource {
	return []*inttypes.ServicePackageSDKResource{
		{
			Factory:  resourceBudgetResourceAssociation,
			TypeName: "aws_servicecatalog_budget_resource_association",
			Name:     "Budget Resource Association",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceConstraint,
			TypeName: "aws_servicecatalog_constraint",
			Name:     "Constraint",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceOrganizationsAccess,
			TypeName: "aws_servicecatalog_organizations_access",
			Name:     "Organizations Access",
			Region:   unique.Make(inttypes.ResourceRegionDisabled()),
		},
		{
			Factory:  resourcePortfolio,
			TypeName: "aws_servicecatalog_portfolio",
			Name:     "Portfolio",
			Tags:     unique.Make(inttypes.ServicePackageResourceTags{}),
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourcePortfolioShare,
			TypeName: "aws_servicecatalog_portfolio_share",
			Name:     "Portfolio Share",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourcePrincipalPortfolioAssociation,
			TypeName: "aws_servicecatalog_principal_portfolio_association",
			Name:     "Principal Portfolio Association",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceProduct,
			TypeName: "aws_servicecatalog_product",
			Name:     "Product",
			Tags:     unique.Make(inttypes.ServicePackageResourceTags{}),
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceProductPortfolioAssociation,
			TypeName: "aws_servicecatalog_product_portfolio_association",
			Name:     "Product Portfolio Association",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceProvisionedProduct,
			TypeName: "aws_servicecatalog_provisioned_product",
			Name:     "Provisioned Product",
			Tags:     unique.Make(inttypes.ServicePackageResourceTags{}),
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceProvisioningArtifact,
			TypeName: "aws_servicecatalog_provisioning_artifact",
			Name:     "Provisioning Artifact",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceServiceAction,
			TypeName: "aws_servicecatalog_service_action",
			Name:     "Service Action",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceTagOption,
			TypeName: "aws_servicecatalog_tag_option",
			Name:     "Tag Option",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
		{
			Factory:  resourceTagOptionResourceAssociation,
			TypeName: "aws_servicecatalog_tag_option_resource_association",
			Name:     "Tag Option Resource Association",
			Region:   unique.Make(inttypes.ResourceRegionDefault()),
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.ServiceCatalog
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*servicecatalog.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))
	optFns := []func(*servicecatalog.Options){
		servicecatalog.WithEndpointResolverV2(newEndpointResolverV2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
		func(o *servicecatalog.Options) {
			if region := config[names.AttrRegion].(string); o.Region != region {
				tflog.Info(ctx, "overriding provider-configured AWS API region", map[string]any{
					"service":         p.ServicePackageName(),
					"original_region": o.Region,
					"override_region": region,
				})
				o.Region = region
			}
		},
		func(o *servicecatalog.Options) {
			if inContext, ok := conns.FromContext(ctx); ok && inContext.VCREnabled() {
				tflog.Info(ctx, "overriding retry behavior to immediately return VCR errors")
				o.Retryer = conns.AddIsErrorRetryables(cfg.Retryer().(aws.RetryerV2), retry.IsErrorRetryableFunc(vcr.InteractionNotFoundRetryableFunc))
			}
		},
		withExtraOptions(ctx, p, config),
	}

	return servicecatalog.NewFromConfig(cfg, optFns...), nil
}

// withExtraOptions returns a functional option that allows this service package to specify extra API client options.
// This option is always called after any generated options.
func withExtraOptions(ctx context.Context, sp conns.ServicePackage, config map[string]any) func(*servicecatalog.Options) {
	if v, ok := sp.(interface {
		withExtraOptions(context.Context, map[string]any) []func(*servicecatalog.Options)
	}); ok {
		optFns := v.withExtraOptions(ctx, config)

		return func(o *servicecatalog.Options) {
			for _, optFn := range optFns {
				optFn(o)
			}
		}
	}

	return func(*servicecatalog.Options) {}
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
