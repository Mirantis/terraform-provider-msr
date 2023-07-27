package provider

import (
	"context"
	"fmt"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &OrgResource{}

type OrgResourceModel struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

type OrgResource struct {
	client client.Client
}

func NewOrgResource() resource.Resource {
	return &OrgResource{}
}

func (r *OrgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org"
}

func (r *OrgResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization",
				Required:            true,
			},
		},
		MarkdownDescription: "Organzation resource",
	}
}

func (r *OrgResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var orgData *OrgResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &orgData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	acc := client.CreateAccount{
		Name:  orgData.Name.ValueString(),
		IsOrg: true,
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr org resource handler is in testing mode, no creation will be run.")
		orgData.Id = basetypes.NewStringValue(TestingVersion)
	} else {
		rAcc, err := r.client.CreateAccount(ctx, acc)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected Create Account error",
				err.Error(),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("created a Org resource `%s`", orgData.Name.ValueString()))

		orgData.Id = basetypes.NewStringValue(rAcc.ID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &orgData)...)
}

func (r *OrgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read org resource")

	var data *OrgResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr org resource handler is in testing mode, no read will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		rAcc, err := r.client.ReadAccount(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		data.Name = types.StringValue(rAcc.Name)
		data.Id = types.StringValue(rAcc.ID)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Can't really update Org resource
	tflog.Trace(ctx, "No action taken. Org resourcs can't be updated.")
}

func (r *OrgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrgResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr org resource handler is in testing mode, no deletion will be run.")
	} else if err := r.client.DeleteAccount(ctx, data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}
}

func (r *OrgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
