package ranger

import (
	"context"

	"github.com/vranyes/goranger/policy"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("RANGER_HOST", nil),
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANGER_USERNAME", "admin"),
			},
			"password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANGER_PASSWORD", nil),
			},
			"skip_ssl_verify": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANGER_SSL", true),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ranger_policy": resourcePolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ranger_policy": dataSourcePolicy(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	skip_ssl_verify, ok := d.Get("skip_ssl_verify").(bool)
	if !ok {
		if d.Get("skip_ssl_verify").(string) == "true" {
			skip_ssl_verify = true
		} else {
			skip_ssl_verify = false
		}
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := policy.NewPolicyClient(host, username, password, skip_ssl_verify)
	return &c, diags
}
