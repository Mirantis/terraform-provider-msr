package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TeamResource{}

type TeamResourceModel struct {
	Name        types.String `tfsdk:"name"`
	OrgID       types.String `tfsdk:"org_id"`
	Description types.String `tfsdk:"description"`
	UserIDs     types.List   `tfsdk:"user_ids"`
	Id          types.String `tfsdk:"id"`
}

type TeamResource struct {
	client client.Client
}

func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

func (r *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"

}

func (r *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				MarkdownDescription: "Name for the team",
				Required:            true,
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "The organization id for the team",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the team",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"user_ids": schema.ListAttribute{
				MarkdownDescription: "The user ids belonging to the team",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
			},
		},
		MarkdownDescription: "Team resource",
	}
}

func (r *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TeamResourceModel

	// // Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	team := client.Team{
		OrgID:       data.OrgID.ValueString(),
		Description: data.Description.ValueString(),
		Name:        data.Name.ValueString(),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr team resource handler is in testing mode, no creation will be run.")
		data.Id = basetypes.NewStringValue(TestingVersion)
	} else {
		rTeam, err := r.client.CreateTeam(ctx, data.OrgID.ValueString(), team)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected Create Team error",
				err.Error(),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("created Team resource `%s`", data.Name.ValueString()))
		data.Id = basetypes.NewStringValue(rTeam.ID)

		var usersSlice []string
		data.UserIDs.ElementsAs(ctx, usersSlice, false)

		for _, id := range usersSlice {
			u := client.ResponseAccount{
				ID: id,
			}
			if err := r.client.AddUserToTeam(ctx, data.OrgID.ValueString(), data.Id.ValueString(), u); err != nil {
				resp.Diagnostics.AddError(
					"Unexpected AddUserToTeam error",
					err.Error(),
				)
				return
			}
			tflog.Trace(ctx, fmt.Sprintf("added user `%s` to team `%s`", id, data.Name.ValueString()))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read team resource")
	var data *TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr team resource handler is in testing mode, no read will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		t, err := r.client.ReadTeam(ctx, data.OrgID.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected ReadTeam error",
				err.Error(),
			)
			return
		}
		data.Id = types.StringValue(t.ID)
		data.Name = types.StringValue(t.Name)
		data.Description = types.StringValue(t.Description)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update team resource")

	var data *TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr team resource handler is in testing mode, no update will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		team := client.Team{
			ID:          data.Id.ValueString(),
			Description: data.Description.ValueString(),
			Name:        data.Name.ValueString(),
		}
		rTeam, err := r.client.UpdateTeam(ctx, data.OrgID.ValueString(), team)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}

		// Overwrite team with refreshed state
		data.Id = types.StringValue(rTeam.ID)
		data.Name = types.StringValue(rTeam.Name)
		data.Description = types.StringValue(rTeam.Description)

		var users []string
		data.UserIDs.ElementsAs(ctx, users, false)
		if err := r.client.UpdateTeamUsers(ctx, data.OrgID.ValueString(), data.Id.ValueString(), users); err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Updated the users of the %s/%s team", data.OrgID, data.Name), map[string]any{"success": true})
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	tflog.Debug(ctx, "Updated 'team' resource", map[string]any{"success": true})
}

func (r *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TeamResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr user resource handler is in testing mode, no deletion will be run.")
	} else if err := r.client.DeleteTeam(ctx, data.OrgID.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted team resource", map[string]any{"success": true})
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: org_id,team_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
