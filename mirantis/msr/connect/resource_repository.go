package connect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mirantis/terraform-provider-msr/mirantis/msr/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceRepo for managing MSR repository
func ResourceRepo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepoCreate,
		ReadContext:   resourceRepoRead,
		UpdateContext: resourceRepoUpdate,
		DeleteContext: resourceRepoDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scan_on_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRepoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	repo := client.CreateRepo{
		Name:       d.Get("name").(string),
		ScanOnPush: d.Get("scan_on_push").(bool),
	}
	orgName := d.Get("org_name").(string)
	_, err := c.CreateRepo(ctx, orgName, repo)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createRepoID(ctx, orgName, repo.Name))

	return diag.Diagnostics{}
}

func resourceRepoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)
	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	orgID, repoID, err := extractRepoIDs(ctx, d.State().ID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = c.ReadRepo(ctx, orgID, repoID)
	if err != nil {
		// If the repo doesn't exist we should gracefully handle it
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(createRepoID(ctx, orgID, repoID))

	return diag.Diagnostics{}
}

func resourceRepoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)

	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}

	orgID, repoID, err := extractRepoIDs(ctx, d.State().ID)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	repo := client.UpdateRepo{
		ScanOnPush: d.Get("scan_on_push").(bool),
		Visibility: "private",
	}

	if d.HasChange("scan_on_push") {
		if _, err := c.UpdateRepo(ctx, orgID, repoID, repo); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRepoRead(ctx, d, m)
}

func resourceRepoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(client.Client)

	if !ok {
		return diag.Errorf("unable to cast meta interface to MSR Client")
	}
	if err := c.DeleteRepo(ctx, d.State().ID); err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diag.Diagnostics{}
}

func extractRepoIDs(ctx context.Context, id string) (org_id string, repo_id string, err error) {
	ids := strings.Split(id, IdDelimiter)

	if len(ids) > 2 || len(ids) < 2 {
		return "", "", fmt.Errorf("resource ID is invalid format '%s'", id)
	}
	return ids[0], ids[1], nil
}

func createRepoID(ctx context.Context, orgID string, teamID string) (id string) {
	return fmt.Sprintf("%s%s%s", orgID, IdDelimiter, teamID)
}
