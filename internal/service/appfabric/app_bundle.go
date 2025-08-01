// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appfabric

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appfabric"
	awstypes "github.com/aws/aws-sdk-go-v2/service/appfabric/types"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	fwflex "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_appfabric_app_bundle", name="App Bundle")
// @Tags(identifierAttribute="id")
// @ArnIdentity(identityDuplicateAttributes="id")
// @Testing(serialize=true)
// @Testing(generator=false)
// @Testing(existsType="github.com/aws/aws-sdk-go-v2/service/appfabric/types;awstypes;awstypes.AppBundle")
// @Testing(preCheckRegion="us-east-1;ap-northeast-1;eu-west-1")
func newAppBundleResource(context.Context) (resource.ResourceWithConfigure, error) {
	r := &appBundleResource{}

	return r, nil
}

type appBundleResource struct {
	framework.ResourceWithModel[appBundleResourceModel]
	framework.WithImportByIdentity
}

func (r *appBundleResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			names.AttrARN: framework.ARNAttributeComputedOnly(),
			"customer_managed_key_arn": schema.StringAttribute{
				CustomType: fwtypes.ARNType,
				Optional:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			names.AttrID:      framework.IDAttribute(),
			names.AttrTags:    tftags.TagsAttribute(),
			names.AttrTagsAll: tftags.TagsAttributeComputedOnly(),
		},
	}
}

func (r *appBundleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data appBundleResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AppFabricClient(ctx)

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		response.Diagnostics.AddError("creating AppFabric App Bundle", err.Error())
	}

	input := &appfabric.CreateAppBundleInput{
		ClientToken:                  aws.String(uuid),
		CustomerManagedKeyIdentifier: fwflex.StringFromFramework(ctx, data.CustomerManagedKeyARN),
		Tags:                         getTagsIn(ctx),
	}

	output, err := conn.CreateAppBundle(ctx, input)

	if err != nil {
		response.Diagnostics.AddError("creating AppFabric App Bundle", err.Error())

		return
	}

	// Set values for unknowns.
	data.ARN = fwflex.StringToFramework(ctx, output.AppBundle.Arn)
	data.setID()

	response.Diagnostics.Append(response.State.Set(ctx, data)...)
}

func (r *appBundleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data appBundleResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AppFabricClient(ctx)

	appBundle, err := findAppBundleByID(ctx, conn, data.ID.ValueString())

	if tfresource.NotFound(err) {
		response.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
		response.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("reading AppFabric App Bundle (%s)", data.ID.ValueString()), err.Error())

		return
	}

	// Set attributes for import.
	response.Diagnostics.Append(fwflex.Flatten(ctx, appBundle, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *appBundleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data appBundleResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().AppFabricClient(ctx)

	input := appfabric.DeleteAppBundleInput{
		AppBundleIdentifier: data.ID.ValueStringPointer(),
	}
	_, err := conn.DeleteAppBundle(ctx, &input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return
	}

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("deleting AppFabric App Bundle (%s)", data.ID.ValueString()), err.Error())

		return
	}
}

func findAppBundleByID(ctx context.Context, conn *appfabric.Client, arn string) (*awstypes.AppBundle, error) {
	input := &appfabric.GetAppBundleInput{
		AppBundleIdentifier: aws.String(arn),
	}

	output, err := conn.GetAppBundle(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.AppBundle == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.AppBundle, nil
}

type appBundleResourceModel struct {
	framework.WithRegionModel
	ARN                   types.String `tfsdk:"arn"`
	CustomerManagedKeyARN fwtypes.ARN  `tfsdk:"customer_managed_key_arn"`
	ID                    types.String `tfsdk:"id"`
	Tags                  tftags.Map   `tfsdk:"tags"`
	TagsAll               tftags.Map   `tfsdk:"tags_all"`
}

func (data *appBundleResourceModel) setID() {
	data.ID = data.ARN
}
