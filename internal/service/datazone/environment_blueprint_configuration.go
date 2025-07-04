// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datazone

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/datazone"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_datazone_environment_blueprint_configuration", name="Environment Blueprint Configuration")
func newEnvironmentBlueprintConfigurationResource(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &environmentBlueprintConfigurationResource{}
	return r, nil
}

const (
	ResNameEnvironmentBlueprintConfiguration = "Environment Blueprint Configuration"
)

type environmentBlueprintConfigurationResource struct {
	framework.ResourceWithModel[environmentBlueprintConfigurationResourceModel]
}

func (r *environmentBlueprintConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled_regions": schema.ListAttribute{
				CustomType:  fwtypes.ListOfStringType,
				ElementType: types.StringType,
				Required:    true,
			},
			"environment_blueprint_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"manage_access_role_arn": schema.StringAttribute{
				CustomType: fwtypes.ARNType,
				Optional:   true,
			},
			"provisioning_role_arn": schema.StringAttribute{
				CustomType: fwtypes.ARNType,
				Optional:   true,
			},
			"regional_parameters": schema.MapAttribute{
				CustomType: fwtypes.MapOfStringType,
				Optional:   true,
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
			},
		},
	}
}

func (r *environmentBlueprintConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Meta().DataZoneClient(ctx)

	var plan environmentBlueprintConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &datazone.PutEnvironmentBlueprintConfigurationInput{
		DomainIdentifier:               plan.DomainId.ValueStringPointer(),
		EnabledRegions:                 flex.ExpandFrameworkStringValueList(ctx, plan.EnabledRegions),
		EnvironmentBlueprintIdentifier: plan.EnvironmentBlueprintId.ValueStringPointer(),
	}

	if !plan.ManageAccessRoleArn.IsNull() {
		in.ManageAccessRoleArn = plan.ManageAccessRoleArn.ValueStringPointer()
	}

	if !plan.ProvisioningRoleArn.IsNull() {
		in.ProvisioningRoleArn = plan.ProvisioningRoleArn.ValueStringPointer()
	}

	if !plan.RegionalParameters.IsNull() {
		var tfMap map[string]map[string]string
		resp.Diagnostics.Append(plan.RegionalParameters.ElementsAs(ctx, &tfMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		in.RegionalParameters = tfMap
	}

	out, err := conn.PutEnvironmentBlueprintConfiguration(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionCreating, ResNameEnvironmentBlueprintConfiguration, plan.EnvironmentBlueprintId.String(), err),
			err.Error(),
		)
		return
	}

	if out == nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionCreating, ResNameEnvironmentBlueprintConfiguration, plan.EnvironmentBlueprintId.String(), nil),
			errors.New("empty output").Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *environmentBlueprintConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Meta().DataZoneClient(ctx)

	var state environmentBlueprintConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := findEnvironmentBlueprintConfigurationByIDs(ctx, conn, state.DomainId.ValueString(), state.EnvironmentBlueprintId.ValueString())
	if tfresource.NotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionSetting, ResNameEnvironmentBlueprintConfiguration, state.EnvironmentBlueprintId.String(), err),
			err.Error(),
		)
		return
	}

	state.DomainId = flex.StringToFramework(ctx, out.DomainId)
	state.EnabledRegions = flex.FlattenFrameworkStringValueListOfStringLegacy(ctx, out.EnabledRegions)
	state.EnvironmentBlueprintId = flex.StringToFramework(ctx, out.EnvironmentBlueprintId)
	state.ManageAccessRoleArn = flex.StringToFrameworkARN(ctx, out.ManageAccessRoleArn)
	state.ProvisioningRoleArn = flex.StringToFrameworkARN(ctx, out.ProvisioningRoleArn)

	regionalParameters, d := flattenRegionalParameters(ctx, &out.RegionalParameters)
	resp.Diagnostics.Append(d...)
	state.RegionalParameters = fwtypes.MapValueOf[types.String]{MapValue: regionalParameters}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentBlueprintConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.Meta().DataZoneClient(ctx)

	var plan, state environmentBlueprintConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.EnabledRegions.Equal(state.EnabledRegions) ||
		!plan.ManageAccessRoleArn.Equal(state.ManageAccessRoleArn) ||
		!plan.ProvisioningRoleArn.Equal(state.ProvisioningRoleArn) ||
		!plan.RegionalParameters.Equal(state.RegionalParameters) {
		in := &datazone.PutEnvironmentBlueprintConfigurationInput{
			DomainIdentifier:               plan.DomainId.ValueStringPointer(),
			EnabledRegions:                 flex.ExpandFrameworkStringValueList(ctx, plan.EnabledRegions),
			EnvironmentBlueprintIdentifier: plan.EnvironmentBlueprintId.ValueStringPointer(),
		}

		if !plan.ManageAccessRoleArn.IsNull() {
			in.ManageAccessRoleArn = plan.ManageAccessRoleArn.ValueStringPointer()
		}

		if !plan.ProvisioningRoleArn.IsNull() {
			in.ProvisioningRoleArn = plan.ProvisioningRoleArn.ValueStringPointer()
		}

		if !plan.RegionalParameters.IsNull() {
			var tfMap map[string]map[string]string
			resp.Diagnostics.Append(plan.RegionalParameters.ElementsAs(ctx, &tfMap, false)...)
			if resp.Diagnostics.HasError() {
				return
			}

			in.RegionalParameters = tfMap
		}

		out, err := conn.PutEnvironmentBlueprintConfiguration(ctx, in)
		if err != nil {
			resp.Diagnostics.AddError(
				create.ProblemStandardMessage(names.DataZone, create.ErrActionUpdating, ResNameEnvironmentBlueprintConfiguration, plan.EnvironmentBlueprintId.String(), err),
				err.Error(),
			)
			return
		}
		if out == nil {
			resp.Diagnostics.AddError(
				create.ProblemStandardMessage(names.DataZone, create.ErrActionUpdating, ResNameEnvironmentBlueprintConfiguration, plan.EnvironmentBlueprintId.String(), nil),
				errors.New("empty output").Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentBlueprintConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	conn := r.Meta().DataZoneClient(ctx)

	var state environmentBlueprintConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &datazone.DeleteEnvironmentBlueprintConfigurationInput{
		DomainIdentifier:               state.DomainId.ValueStringPointer(),
		EnvironmentBlueprintIdentifier: state.EnvironmentBlueprintId.ValueStringPointer(),
	}

	_, err := conn.DeleteEnvironmentBlueprintConfiguration(ctx, in)
	if err != nil {
		if isResourceMissing(err) {
			return
		}
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.DataZone, create.ErrActionDeleting, ResNameEnvironmentBlueprintConfiguration, state.EnvironmentBlueprintId.String(), err),
			err.Error(),
		)
		return
	}
}

func (r *environmentBlueprintConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Resource Import Invalid ID", fmt.Sprintf("Wrong format for import ID (%s), use: 'domain-id/environment-blueprint-id'", req.ID))
		return
	}
	domainId := parts[0]
	environmentBlueprintId := parts[1]

	environmentBlueprintConfiguration, err := findEnvironmentBlueprintConfigurationByIDs(ctx, r.Meta().DataZoneClient(ctx), domainId, environmentBlueprintId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Importing Resource",
			err.Error(),
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), aws.ToString(environmentBlueprintConfiguration.DomainId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_blueprint_id"), aws.ToString(environmentBlueprintConfiguration.EnvironmentBlueprintId))...)
}

func findEnvironmentBlueprintConfigurationByIDs(ctx context.Context, conn *datazone.Client, domainId, environmentBlueprintId string) (*datazone.GetEnvironmentBlueprintConfigurationOutput, error) {
	in := &datazone.GetEnvironmentBlueprintConfigurationInput{
		DomainIdentifier:               aws.String(domainId),
		EnvironmentBlueprintIdentifier: aws.String(environmentBlueprintId),
	}

	out, err := conn.GetEnvironmentBlueprintConfiguration(ctx, in)
	if err != nil {
		if isResourceMissing(err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func flattenRegionalParameters(ctx context.Context, apiObject *map[string]map[string]string) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := types.MapType{ElemType: types.StringType}

	if apiObject == nil || len(*apiObject) == 0 {
		return types.MapNull(elemType), diags
	}

	elements := map[string]types.Map{}

	for k, v := range *apiObject {
		elements[k] = flex.FlattenFrameworkStringValueMap(ctx, v)
	}

	mapVal, d := types.MapValueFrom(ctx, types.MapType{ElemType: types.StringType}, elements)
	diags.Append(d...)

	return mapVal, diags
}

type environmentBlueprintConfigurationResourceModel struct {
	framework.WithRegionModel
	DomainId               types.String         `tfsdk:"domain_id"`
	EnabledRegions         fwtypes.ListOfString `tfsdk:"enabled_regions"`
	EnvironmentBlueprintId types.String         `tfsdk:"environment_blueprint_id"`
	ManageAccessRoleArn    fwtypes.ARN          `tfsdk:"manage_access_role_arn"`
	ProvisioningRoleArn    fwtypes.ARN          `tfsdk:"provisioning_role_arn"`
	RegionalParameters     fwtypes.MapOfString  `tfsdk:"regional_parameters"`
}
