package provider

import (
	"context"
	"fmt"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &accountDataSource{}
)

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSource struct {
	client client.Client
}

// accountDataSourceModel maps the data source schema data.
type accountDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	NameOrID     types.String `tfsdk:"name_or_id"`
	Name         types.String `tfsdk:"name"`
	FullName     types.String `tfsdk:"full_name"`
	IsOrg        types.Bool   `tfsdk:"is_org"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	IsAdmin      types.Bool   `tfsdk:"is_admin"`
	IsImported   types.Bool   `tfsdk:"is_imported"`
	OnDemand     types.Bool   `tfsdk:"on_demand"`
	OtpEnabled   types.Bool   `tfsdk:"otp_enabled"`
	MembersCount types.Int64  `tfsdk:"members_count"`
	TeamsCount   types.Int64  `tfsdk:"teams_count"`
}

// Configure adds the provider configured client to the data source.
func (d *accountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)
	if !ok {
		tflog.Error(ctx, "Unable to prepare client")
		return
	}
	d.client = client

}

func (d *accountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *accountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
			},
			"name_or_id": schema.StringAttribute{
				MarkdownDescription: "The name or id of the account",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the account",
				Optional:            true,
				Computed:            true,
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "The full name of the account",
				Optional:            true,
				Computed:            true,
			},
			"is_org": schema.BoolAttribute{
				MarkdownDescription: "Is the account organization",
				Optional:            true,
				Computed:            true,
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Is the account admin",
				Optional:            true,
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Is the account active",
				Optional:            true,
				Computed:            true,
			},
			"is_imported": schema.BoolAttribute{
				MarkdownDescription: "Is the account imported",
				Optional:            true,
				Computed:            true,
			},
			"on_demand": schema.BoolAttribute{
				MarkdownDescription: "Is the account on demand",
				Optional:            true,
				Computed:            true,
			},
			"otp_enabled": schema.BoolAttribute{
				MarkdownDescription: "Is `otp_enabled` for the account",
				Optional:            true,
				Computed:            true,
			},
			"members_count": schema.Int64Attribute{
				MarkdownDescription: "The members count if the the account is org",
				Optional:            true,
				Computed:            true,
			},
			"teams_count": schema.Int64Attribute{
				MarkdownDescription: "The number of teams for the account",
				Optional:            true,
				Computed:            true,
			},
		},
		MarkdownDescription: "Fetch an account",
	}
}

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read account data source")
	var data accountDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if d.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr account datasource handler is in testing mode, no injestion will be run.")
		data.ID = basetypes.NewStringValue(TestingVersion)
	} else {
		rAcc, err := d.client.ReadAccount(ctx, data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Item",
				err.Error(),
			)
			return
		}
		data.ID = types.StringValue(rAcc.ID)
		data.Name = types.StringValue(rAcc.Name)
		data.FullName = types.StringValue(rAcc.FullName)
		data.IsOrg = types.BoolValue(rAcc.IsOrg)
		data.IsActive = types.BoolValue(rAcc.IsActive)
		data.IsAdmin = types.BoolValue(rAcc.IsAdmin)
		data.IsImported = types.BoolValue(rAcc.IsImported)
		data.OnDemand = types.BoolValue(rAcc.OnDemand)
		data.OtpEnabled = types.BoolValue(rAcc.OtpEnabled)
		data.MembersCount = types.Int64Value(int64(rAcc.MembersCount))
		data.TeamsCount = types.Int64Value(int64(rAcc.TeamsCount))

		tflog.Trace(ctx, fmt.Sprintf("read in account data source `%s`", rAcc.ID))
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "Finished reading account data source", map[string]any{"success": true})
}
