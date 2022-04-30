package ranger

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vranyes/goranger/policy"
)

func dataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*policy.PolicyClient)

	id := d.Get("id").(int)
	policy, err := client.GetPolicyById(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", policy.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return diags
}
