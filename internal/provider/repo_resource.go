package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &RepoResource{}

type RepoResourceModel struct {
	Name       types.String `tfsdk:"name"`
	OrgName    types.String `tfsdk:"org_name"`
	ScanOnPush types.Bool   `tfsdk:"scan_on_push"`
	Visibility types.String `tfsdk:"visibility"`
	Id         types.String `tfsdk:"id"`
}

type RepoResource struct {
	client client.Client
}

func NewRepoResource() resource.Resource {
	return &RepoResource{}
}

func (r *RepoResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repo"

}

func (r *RepoResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				MarkdownDescription: "The name of the team",
				Required:            true,
			},
			"org_name": schema.StringAttribute{
				MarkdownDescription: "The organization name for the repo",
				Required:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "The visibility of the the repo",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("private"),
			},
			"scan_on_push": schema.BoolAttribute{
				MarkdownDescription: "The scan ",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
		MarkdownDescription: "Repo resource",
	}
}

func (r *RepoResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RepoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RepoResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	repo := client.CreateRepo{
		Name:       data.Name.ValueString(),
		ScanOnPush: data.ScanOnPush.ValueBool(),
		Visibility: data.Visibility.ValueString(),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr repo resource handler is in testing mode, no creation will be run.")
		data.Id = basetypes.NewStringValue(TestingVersion)
	} else {
		rRepo, err := r.client.CreateRepo(ctx, data.OrgName.ValueString(), repo)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected Create Team error",
				err.Error(),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("created Org resource `%s`", data.Name.ValueString()))
		data.Id = basetypes.NewStringValue(rRepo.ID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read repo resource")
	var data *RepoResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr repo resource handler is in testing mode, no read will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		r, err := r.client.ReadRepo(ctx, data.OrgName.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected ReadTeam error",
				err.Error(),
			)
			return
		}
		data.Id = types.StringValue(r.ID)
		data.Name = types.StringValue(r.Name)
		data.ScanOnPush = types.BoolValue(r.ScanOnPush)
		data.Visibility = types.StringValue(r.Visibility)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RepoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update team resource")

	var data *RepoResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr repo resource handler is in testing mode, no update will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		repo := client.UpdateRepo{
			ScanOnPush: data.ScanOnPush.ValueBool(),
			Visibility: data.Visibility.ValueString(),
		}
		rRepo, err := r.client.UpdateRepo(ctx, data.OrgName.ValueString(), data.Name.ValueString(), repo)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}

		data.ScanOnPush = types.BoolValue(rRepo.ScanOnPush)
		data.Id = types.StringValue(rRepo.ID)
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	tflog.Debug(ctx, "Updated 'repo' resource", map[string]any{"success": true})
}

func (r *RepoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RepoResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr user resource handler is in testing mode, no deletion will be run.")
	} else if err := r.client.DeleteRepo(ctx, data.OrgName.ValueString(), data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted team resource", map[string]any{"success": true})
}

func (r *RepoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: org_name,repo_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
