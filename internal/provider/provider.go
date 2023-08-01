package provider

import (
	"context"

	//	"net/http"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	TestingVersion = "test"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &MSRProvider{}

// MSRProvider defines the provider implementation.
type MSRProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MSRProviderModel describes the provider data model.
type MSRProviderModel struct {
	Host            types.String `tfsdk:"host"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	UnsafeSSLClient types.Bool   `tfsdk:"unsafe_ssl_client"`
}

func (p *MSRProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "msr"
	resp.Version = p.version
}

func (p *MSRProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The host url of the MSR instance",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username to login in the MSR instance",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to login in the MSR instance",
				Required:            true,
				Sensitive:           true,
			},
			"unsafe_ssl_client": schema.BoolAttribute{
				MarkdownDescription: "Use of unsafe SSL client",
				Optional:            true,
			},
		},
	}
}

func (p *MSRProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MSRProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	testMode := false
	if p.version == TestingVersion {
		testMode = true
	}

	var c client.Client
	var err error
	if data.UnsafeSSLClient.ValueBool() {
		c, err = client.NewUnsafeSSLClient(data.Host.ValueString(), data.Username.ValueString(), data.Password.ValueString(), testMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create NewUnsafeSSLClient from terraform config",
				err.Error(),
			)

			return
		}
	} else {
		c, err = client.NewDefaultClient(data.Host.ValueString(), data.Username.ValueString(), data.Password.ValueString(), testMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create NewDefaultClient from terraform config",
				err.Error(),
			)

			return
		}
	}

	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *MSRProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrgResource,
		NewUserResource,
	}
}

func (p *MSRProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Account
		// Accounts
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MSRProvider{
			version: version,
		}
	}
}
