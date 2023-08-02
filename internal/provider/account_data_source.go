// import (
// 	"context"
// 	"strconv"
// 	"time"

// 	"github.com/Mirantis/terraform-provider-msr/internal/client"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// // DataSourceAccount for retrieving a single MSR
// func dataSourceAccount() *schema.Resource {
// 	return &schema.Resource{
// 		ReadContext: dataSourceAccountRead,
// 		Schema: map[string]*schema.Schema{
// 			"name_or_id": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"name": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"full_name": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"id": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"is_org": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"is_admin": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"is_active": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"is_imported": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"on_demand": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"otp_enabled": {
// 				Type:     schema.TypeBool,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"members_count": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 				Computed: true,
// 			},
// 			"teams_count": {
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 				Computed: true,
// 			},
// 		},
// 	}
// }

// func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c, ok := m.(client.Client)
// 	if !ok {
// 		return diag.Errorf("unable to cast meta interface to MSR Client")
// 	}

// 	accNameOrID := d.Get("name_or_id").(string)
// 	rAccount, err := c.ReadAccount(ctx, accNameOrID)
// 	if err != nil {
// 		// If the accounts doesn't exist we should gracefully handle it
// 		d.SetId("")
// 		return diag.FromErr(err)
// 	}

// 	// Common fields between user and org
// 	if err := d.Set("name", rAccount.Name); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("id", rAccount.ID); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("is_org", rAccount.IsOrg); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("full_name", rAccount.FullName); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	// Org specific fields
// 	if rAccount.IsOrg {
// 		if err := d.Set("members_count", rAccount.MembersCount); err != nil {
// 			return diag.FromErr(err)
// 		}
// 		if err := d.Set("teams_count", rAccount.TeamsCount); err != nil {
// 			return diag.FromErr(err)
// 		}
// 	} else { // User specific fields
// 		if err := d.Set("is_admin", rAccount.IsAdmin); err != nil {
// 			return diag.FromErr(err)
// 		}
// 		if err := d.Set("is_active", rAccount.IsActive); err != nil {
// 			return diag.FromErr(err)
// 		}
// 		if err := d.Set("is_imported", rAccount.IsImported); err != nil {
// 			return diag.FromErr(err)
// 		}
// 		if err := d.Set("on_demand", rAccount.OnDemand); err != nil {
// 			return diag.FromErr(err)
// 		}
// 		if err := d.Set("otp_enabled", rAccount.OtpEnabled); err != nil {
// 			return diag.FromErr(err)
// 		}
// 	}

// 	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

// 	return diag.Diagnostics{}
// }

package provider

// import (
// 	"context"
// 	"fmt"

// 	"github.com/Mirantis/terraform-provider-msr/internal/client"
// 	"github.com/hashicorp/terraform-plugin-framework/path"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// )

// var _ resource.Resource = &OrgResource{}

// type OrgResourceModel struct {
// 	Name types.String `tfsdk:"name"`
// 	Id   types.String `tfsdk:"id"`
// }

// type OrgResource struct {
// 	client client.Client
// }

// func NewOrgResource() resource.Resource {
// 	return &OrgResource{}
// }

// func (r *OrgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_account"
// }

// 	// 			"name": {
// 	// 				Type:     schema.TypeString,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"full_name": {
// 	// 				Type:     schema.TypeString,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"id": {
// 	// 				Type:     schema.TypeString,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"is_org": {
// 	// 				Type:     schema.TypeBool,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"is_admin": {
// 	// 				Type:     schema.TypeBool,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"is_active": {
// 	// 				Type:     schema.TypeBool,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"is_imported": {
// 	// 				Type:     schema.TypeBool,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"on_demand": {
// 	// 				Type:     schema.TypeBool,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"otp_enabled": {
// 	// 				Type:     schema.TypeBool,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"members_count": {
// 	// 				Type:     schema.TypeInt,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},
// 	// 			"teams_count": {
// 	// 				Type:     schema.TypeInt,
// 	// 				Optional: true,
// 	// 				Computed: true,
// 	// 			},

// func (r *OrgResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		// This description is used by the documentation generator and the language server.

// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				Computed:            true,
// 				MarkdownDescription: "Identifier",
// 				PlanModifiers: []planmodifier.String{
// 					stringplanmodifier.UseStateForUnknown(),
// 				},
// 			},
// 			"name_or_id": schema.StringAttribute{
// 				MarkdownDescription: "The name or id of the account",
// 				Required:            true,
// 			},
// 			"name": schema.StringAttribute{
// 				MarkdownDescription: "The name or id of the account",
// 				Optional:            true,
// 				Computed: true,
// 				Default: stringdefault.StaticString("")
// 			},
// 		},
// 		MarkdownDescription: "Organzation resource",
// 	}
// }

// func (r *OrgResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	// Prevent panic if the provider has not been configured.
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	client, ok := req.ProviderData.(client.Client)
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Client error",
// 			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)

// 		return
// 	}

// 	r.client = client
// }

// func (r *OrgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	tflog.Debug(ctx, "Preparing to read org resource")

// 	var data *OrgResourceModel

// 	// Read Terraform prior state data into the model
// 	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	if r.client.TestMode {
// 		resp.Diagnostics.AddWarning("testing mode warning", "msr org resource handler is in testing mode, no read will be run.")
// 		data.Id = types.StringValue(TestingVersion)
// 	} else {
// 		rAcc, err := r.client.ReadAccount(ctx, data.Name.ValueString())
// 		if err != nil {
// 			resp.Diagnostics.AddError("Client Error", err.Error())
// 			return
// 		}
// 		data.Name = types.StringValue(rAcc.Name)
// 		data.Id = types.StringValue(rAcc.ID)
// 	}

// 	// Save updated data into Terraform state
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }
