package prismatic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/graphql"
	"strconv"
	"time"
)

func dataSourceIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntegrationsRead,
		Schema: map[string]*schema.Schema{
			"integrations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"integration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"integration_definition": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIntegrationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*graphql.Client)

	var diags diag.Diagnostics

	var query struct {
		Integrations struct {
			Nodes []struct {
				Id         string
				Name       string
				Definition string
			}
		}
	}

	if err := client.Query(context.Background(), &query, nil); err != nil {
		return diag.FromErr(err)
	}

	count := len(query.Integrations.Nodes)
	integrations := make([]interface{}, count, count)
	for i, integrationNode := range query.Integrations.Nodes {
		integration := make(map[string]interface{})
		integration["integration_id"] = integrationNode.Id
		integration["integration_name"] = integrationNode.Name
		integration["integration_definition"] = integrationNode.Definition
		integrations[i] = integration
	}

	if err := d.Set("integrations", integrations); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
