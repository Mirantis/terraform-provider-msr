package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PruningPolicyResource{}

type PruningPolicyResourceModel struct {
	Id       types.String                    `tfsdk:"id"`
	Enabled  types.Bool                      `tfsdk:"enabled"`
	OrgName  types.String                    `tfsdk:"org_name"`
	RepoName types.String                    `tfsdk:"repo_name"`
	Rules    []client.PruningPolicyRuleTFSDK `tfsdk:"rule"`
}

type PruningPolicyResource struct {
	client client.Client
}

func NewPruningPolicyResource() resource.Resource {
	return &PruningPolicyResource{}
}

func (r *PruningPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pruning_policy"

}

func (r *PruningPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Pruning policy resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Is the pruning policy enabled",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"org_name": schema.StringAttribute{
				MarkdownDescription: "The organization that contains the repo",
				Required:            true,
			},
			"repo_name": schema.StringAttribute{
				MarkdownDescription: "The repository to apply the pruning policy on",
				Required:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				MarkdownDescription: "The rules of the pruning policy",
				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"field": schema.StringAttribute{
							MarkdownDescription: "The field for the rule",
							Required:            true,
						},
						"operator": schema.StringAttribute{
							MarkdownDescription: "The operator for the particular field",
							Required:            true,
						},
						"values": schema.ListAttribute{
							MarkdownDescription: "The regex values for the rule",
							Required:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *PruningPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)
	client.HTTPClient.Timeout = 240 * time.Second
	if !ok {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PruningPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Preparing to create pruning policy resource")
	var data PruningPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	pruningPolicy := client.CreatePruningPolicy{
		Enabled: true,
		Rules:   client.PruningPolicyRulesToAPI(ctx, data.Rules),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr repo resource handler is in testing mode, no creation will be run.")
		data.Id = basetypes.NewStringValue(TestingVersion)
	} else {
		existingPolicies, err := r.client.ReadPruningPolicies(ctx, data.OrgName.ValueString(), data.RepoName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected Create pruning policy error",
				err.Error(),
			)
			return
		}
		existingPolicy := r.client.PruningPolicyExists(ctx, pruningPolicy, existingPolicies)
		tflog.Debug(ctx, fmt.Sprintf("Returned existingPolicy %+v\n", existingPolicy))

		// There is an existing pruning policy
		if existingPolicy.ID != "" {
			resp.Diagnostics.AddError(
				"Cannot create duplicate pruning policy",
				fmt.Sprintf("Pruning policy already exists with id %s", existingPolicy.ID),
			)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Proceeding with Policy creation %+v\n", existingPolicy))

		rPolicy, err := r.client.CreatePruningPolicy(ctx, data.OrgName.ValueString(), data.RepoName.ValueString(), pruningPolicy)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected Create pruning policy error",
				err.Error(),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("created Pruning policy resource with ID `%s`", data.Id.ValueString()))
		data.Id = basetypes.NewStringValue(rPolicy.ID)
		data.Rules = client.PruningPolicyRulesToTFSDK(ctx, rPolicy.Rules)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PruningPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read pruning policy resource")
	var data PruningPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr pruning policy resource handler is in testing mode, no read will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		rPolicy, err := r.client.ReadPruningPolicy(ctx, data.OrgName.ValueString(), data.RepoName.ValueString(), data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		data.Id = types.StringValue(rPolicy.ID)
		data.Enabled = basetypes.NewBoolValue(rPolicy.Enabled)
		data.Rules = client.PruningPolicyRulesToTFSDK(ctx, rPolicy.Rules)

	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PruningPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update pruning policy resource")

	var data PruningPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr pruning policy resource handler is in testing mode, no update will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		policy := client.CreatePruningPolicy{
			Enabled: data.Enabled.ValueBool(),
			Rules:   client.PruningPolicyRulesToAPI(ctx, data.Rules),
		}
		rPolicy, err := r.client.UpdatePruningPolicy(ctx, data.OrgName.ValueString(), data.RepoName.ValueString(), policy, data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}

		// Overwrite pruning policy with refreshed state
		data.Id = types.StringValue(rPolicy.ID)
		data.Enabled = types.BoolValue(rPolicy.Enabled)
		data.Rules = client.PruningPolicyRulesToTFSDK(ctx, rPolicy.Rules)
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	// Set refreshed state
	tflog.Debug(ctx, fmt.Sprintf("Updated Pruning Policy with ID %s of for %s/%s ", data.Id, data.OrgName, data.RepoName), map[string]any{"success": true})
}

func (r *PruningPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete pruning policy resource")
	var data *PruningPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr pruning policy resource handler is in testing mode, no deletion will be run.")
	} else if err := r.client.DeletePruningPolicy(ctx, data.OrgName.ValueString(), data.RepoName.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted pruning policy resource", map[string]any{"success": true})
}

func (r *PruningPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: org_name,repo_name,pruning_policy_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repo_name"), idParts[1])...)
	// policy ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[2])...)
}
