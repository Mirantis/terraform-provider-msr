// func dataSourceAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c, ok := m.(client.Client)
// 	if !ok {
// 		return diag.Errorf("unable to cast meta interface to MSR Client")
// 	}

// 	filter := client.AccountFilter("all")
// 	if _, ok := d.GetOk("filter"); ok {
// 		inputFilter := d.Get("filter").(string)
// 		filter = client.AccountFilter(inputFilter)

// 		if filter.APIFormOfFilter() != inputFilter {
// 			d.SetId("")
// 			return diag.FromErr(fmt.Errorf("%w. Filter '%s'", client.ErrInvalidFilter, inputFilter))
// 		}
// 	}
// 	rAccounts, err := c.ReadAccounts(ctx, filter)
// 	if err != nil {
// 		// If the accounts doesn't exist we should gracefully handle it
// 		d.SetId("")
// 		return diag.FromErr(err)
// 	}

// 	accounts := make([]map[string]interface{}, 0, len(rAccounts))

// 	for _, u := range rAccounts {
// 		accounts = append(accounts, map[string]interface{}{
// 			"id":            u.ID,
// 			"name":          u.Name,
// 			"full_name":     u.FullName,
// 			"is_active":     u.IsActive,
// 			"is_admin":      u.IsAdmin,
// 			"is_org":        u.IsOrg,
// 			"members_count": u.MembersCount,
// 			"teams_count":   u.TeamsCount,
// 		})
// 	}

// 	if err := d.Set("accounts", accounts); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

// 	return diag.Diagnostics{}
// }

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &accountsDataSource{}
)

func NewaccountsDataSource() datasource.DataSource {
	return &accountsDataSource{}
}

type accountsDataSource struct {
	client client.Client
}

// accountsDataSourceModel maps the data source schema data.
type accountsDataSourceModel struct {
	ID       types.String             `tfsdk:"id"`
	Filter   types.String             `tfsdk:"filter"`
	Accounts []accountDataSourceModel `tfsdk:"accounts"`
}

// Configure adds the provider configured client to the data source.
func (d *accountsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
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

func (d *accountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

func (d *accountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Accounts data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
			},
			"filter": schema.StringAttribute{
				MarkdownDescription: "The name of the account",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.OneOf("users", "orgs", "admins", "non-admins", "active-users", "all")},
			},
		},

		Blocks: map[string]schema.Block{
			"accounts": schema.ListNestedBlock{
				MarkdownDescription: "The accounts retrieved from MSR",
				NestedObject: schema.NestedBlockObject{

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
				},
			},
		},
	}
}

func (d *accountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read account data source")
	var data accountsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if d.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr account datasource handler is in testing mode, no injestion will be run.")
		data.ID = basetypes.NewStringValue(TestingVersion)
	} else {
		filter := client.AccountFilter(data.Filter.ValueString())

		rAccs, err := d.client.ReadAccounts(ctx, filter)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Accounts",
				err.Error(),
			)
			return
		}

		var accs []accountDataSourceModel

		for _, u := range rAccs {
			acc := accountDataSourceModel{
				ID:           basetypes.NewStringValue(u.ID),
				NameOrID:     basetypes.NewStringValue(u.ID),
				Name:         basetypes.NewStringValue(u.Name),
				FullName:     basetypes.NewStringValue(u.FullName),
				IsActive:     basetypes.NewBoolValue(u.IsActive),
				IsAdmin:      basetypes.NewBoolValue(u.IsAdmin),
				IsOrg:        basetypes.NewBoolValue(u.IsOrg),
				MembersCount: basetypes.NewInt64Value(int64(u.MembersCount)),
				TeamsCount:   basetypes.NewInt64Value(int64(u.TeamsCount)),
			}
			accs = append(accs, acc)
		}

		data.Accounts = accs
		data.ID = basetypes.NewStringValue(time.Now().Format(time.RFC850))

		tflog.Trace(ctx, fmt.Sprintf("read in accounts data source `%s`", data.ID))
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "Finished reading accounts data source", map[string]any{"success": true})
}
