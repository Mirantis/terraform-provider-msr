package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/Mirantis/terraform-provider-msr/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceAccount for retrieving a single MSR
func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"name_or_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_org": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_imported": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"on_demand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"otp_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"members_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"teams_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	accNameOrID := d.Get("name_or_id").(string)
	rAccount, err := c.ReadAccount(ctx, accNameOrID)
	if err != nil {
		// If the accounts doesn't exist we should gracefully handle it
		d.SetId("")
		return diag.FromErr(err)
	}

	// Common fields between user and org
	if err := d.Set("name", rAccount.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("id", rAccount.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_org", rAccount.IsOrg); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("full_name", rAccount.FullName); err != nil {
		return diag.FromErr(err)
	}

	// Org specific fields
	if rAccount.IsOrg {
		if err := d.Set("members_count", rAccount.MembersCount); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("teams_count", rAccount.TeamsCount); err != nil {
			return diag.FromErr(err)
		}
	} else { // User specific fields
		if err := d.Set("is_admin", rAccount.IsAdmin); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("is_active", rAccount.IsActive); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("is_imported", rAccount.IsImported); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("on_demand", rAccount.OnDemand); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("otp_enabled", rAccount.OtpEnabled); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}
