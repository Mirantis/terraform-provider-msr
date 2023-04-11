package connect

import (
	"context"

	"github.com/Mirantis/terraform-provider-msr/mirantis/msr/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSR_HOST_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MSR_ADMIN_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MSR_ADMIN_PASS", nil),
			},
			"unsafe_ssl_client": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("MSR_UNSAFE_CLIENT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"msr_user": ResourceUser(),
			"msr_org":  ResourceOrg(),
			"msr_team": ResourceTeam(),
			"msr_repo": ResourceRepo(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"msr_accounts": dataSourceAccounts(),
			"msr_account":  dataSourceAccount(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)
	unsafeClient := d.Get("unsafe_ssl_client").(bool)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error
	var c client.Client
	if unsafeClient {
		c, err = client.NewUnsafeSSLClient(host, username, password)

	} else {
		c, err = client.NewDefaultClient(host, username, password)
	}
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create MSR client",
			Detail:   err.Error(),
		})

		return nil, diags
	}

	healthy, err := c.IsHealthy(ctx)
	if !healthy {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "MSR endpoint is not healthy",
			Detail:   err.Error(),
		})
		return nil, diags
	}
	return c, diags
}
