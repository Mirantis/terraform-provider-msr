package provider

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"time"

// 	"github.com/Mirantis/terraform-provider-msr/internal/client"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// // DataSourceAccounts for retrieving MSR accounts in bulk
// func dataSourceAccounts() *schema.Resource {
// 	return &schema.Resource{
// 		ReadContext: dataSourceAccountsRead,
// 		Schema: map[string]*schema.Schema{
// 			"filter": {
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"accounts": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"id": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"name": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"full_name": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"is_active": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
// 						},
// 						"is_admin": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
// 						},
// 						"is_org": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
// 						},
// 						"members_count": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"teams_count": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

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
