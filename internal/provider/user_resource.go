package provider

import (
	"context"
	"fmt"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &UserResource{}

type UserResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
	FullName types.String `tfsdk:"full_name"`
	IsAdmin  types.Bool   `tfsdk:"is_admin"`
	Id       types.String `tfsdk:"id"`
}

type UserResource struct {
	client client.Client
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				MarkdownDescription: "The name of the user",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthBetween(3, 32)},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password of the user",
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators:          []validator.String{stringvalidator.LengthBetween(8, 16)},
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "The full name of the user",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Is the user an admin",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Optional:            true,
			},
		},
		MarkdownDescription: "User resource",
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	pass := data.Password.ValueString()
	if pass == "" {
		pass = client.GeneratePass()
		data.Password = basetypes.NewStringValue(pass)
	}

	acc := client.CreateAccount{
		Name:       data.Name.ValueString(),
		Password:   pass,
		FullName:   data.FullName.ValueString(),
		IsAdmin:    data.IsAdmin.ValueBool(),
		IsOrg:      false,
		SearchLDAP: false,
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr user resource handler is in testing mode, no creation will be run.")
		data.Id = basetypes.NewStringValue(TestingVersion)
	} else {
		rAcc, err := r.client.CreateAccount(ctx, acc)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected Create Account error",
				err.Error(),
			)
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("created User resource `%s`", data.Name.ValueString()))

		data.Id = basetypes.NewStringValue(rAcc.ID)
		data.Password = basetypes.NewStringValue(pass)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read user resource")
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr user resource handler is in testing mode, no read will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		rAcc, err := r.client.ReadAccount(ctx, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}
		data.Id = types.StringValue(rAcc.ID)
		data.Name = types.StringValue(rAcc.Name)
		data.FullName = types.StringValue(rAcc.FullName)
		data.IsAdmin = types.BoolValue(rAcc.IsAdmin)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update user resource")

	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr user resource handler is in testing mode, no update will be run.")
		data.Id = types.StringValue(TestingVersion)
	} else {
		user := client.UpdateAccount{
			FullName: data.FullName.ValueString(),
			IsAdmin:  data.IsAdmin.ValueBool(),
		}
		rAcc, err := r.client.UpdateAccount(ctx, data.Id.ValueString(), user)
		tflog.Debug(ctx, fmt.Sprintf("The retuerned 'user' %+v", rAcc))

		if err != nil {
			resp.Diagnostics.AddError("Client Error", err.Error())
			return
		}

		// Overwrite user with refreshed state
		data.Id = types.StringValue(rAcc.ID)
		data.Name = types.StringValue(rAcc.Name)
		data.FullName = types.StringValue(rAcc.FullName)
		data.IsAdmin = types.BoolValue(rAcc.IsAdmin)
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	tflog.Debug(ctx, "Updated 'user' resource", map[string]any{"success": true})
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if r.client.TestMode {
		resp.Diagnostics.AddWarning("testing mode warning", "msr user resource handler is in testing mode, no deletion will be run.")
	} else if err := r.client.DeleteAccount(ctx, data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleted user resource", map[string]any{"success": true})
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
