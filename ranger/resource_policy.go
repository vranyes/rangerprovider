package ranger

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vranyes/goranger/policy"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Delete: resourcePolicyDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type: schema.TypeString,
			},
			"name": {
				Type: schema.TypeString,
			},
			"id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: false,
			},
			"enabled": {
				Type:    schema.TypeBool,
				Default: true,
			},
			"labels": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policies": {
				Type: schema.TypeList,
				Elem: map[string]*schema.Schema{
					"accesses": {
						Type: schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
					"users": {
						Type: schema.TypeList,
					},
					"groups": {
						Type: schema.TypeList,
					},
				},
			},
			"policy_type": {
				Type: schema.TypeInt,
			},
			"resources": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: map[string]*schema.Schema{
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"recursive": {
							Type: schema.TypeBool,
						},
						"excludes": {
							Type: schema.TypeBool,
						},
					},
				},
			},
			"service": {
				Type: schema.TypeString,
			},
			"service_type": {
				Type: schema.TypeString,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func policyResourceExpander(input map[string]interface{}) map[string]policy.PolicyResource {
	var prExpanded map[string]policy.PolicyResource
	for k, v := range input {
		remapped := v.(map[string]interface{})
		prExpanded[k] = policy.PolicyResource{
			IsExcludes:  remapped["excludes"].(bool),
			IsRecursive: remapped["recursive"].(bool),
			Values:      remapped["values"].([]string),
		}
	}
	return prExpanded
}

// expands an array of PolicyItemAccess
func accessExpandeder(input []interface{}) []policy.PolicyItemAccess {
	var pia []policy.PolicyItemAccess

	// each v is a PolicyItemAccess
	for _, v := range input {
		remapped := v.(map[string]interface{})
		pia = append(pia, policy.PolicyItemAccess{
			Type:      remapped["type"].(string),
			IsAllowed: remapped["allowed"].(bool),
		})
	}

	return pia
}

// expands a list of PolicyItems
func policyPolicyExpander(input []interface{}) []policy.PolicyItem {
	var pis []policy.PolicyItem

	for _, v := range input {
		p := v.(map[string]interface{})

		pis = append(pis, policy.PolicyItem{
			Accesses: accessExpandeder(p["accesses"].([]interface{})),
			Groups:   p["groups"].([]string),
			Users:    p["users"].([]string),
		})
	}

	return pis
}

func resourcePolicyCreate(d *schema.ResourceData, m interface{}) error {
	var err error
	client := m.(*policy.PolicyClient)

	p := policy.Policy{}

	p.Resources = policyResourceExpander(d.Get("resources").(map[string]interface{}))
	p.PolicyItems = policyPolicyExpander(d.Get("policies").([]interface{}))
	p.PolicyLabels = d.Get("labels").([]string)
	p.Description = d.Get("description").(string)
	p.Name = d.Get("name").(string)
	p.IsEnabled = d.Get("enabled").(bool)
	p.PolicyType = d.Get("policy_type").(int)
	p.Service = d.Get("service").(string)
	p.ServiceType = d.Get("service_type").(string)
	p.CreateTime = d.Get("create_time").(int64)
	p.CreatedBy = d.Get("created_by").(string)
	p.UpdateTime = d.Get("update_time").(int64)
	p.UpdatedBy = d.Get("updated_by").(string)
	p.Guid = d.Get("guid").(string)

	createdPolicy, err := client.CreatePolicy(p)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(createdPolicy.Id))

	return resourcePolicyRead(d, m)
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	var err error
	client := m.(*policy.PolicyClient)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	p, err := client.GetPolicyById(id)
	if err != nil {
		// if err'd on getting, del from state
		d.SetId("")
		return err
	}

	d.Set("resources", p.Resources)
	d.Set("policies", p.PolicyItems)
	d.Set("labels", p.PolicyLabels)
	d.Set("description", p.Description)
	d.Set("name", p.Name)
	d.Set("enabled", p.IsEnabled)
	d.Set("policy_type", p.PolicyType)
	d.Set("service", p.Service)
	d.Set("service_type", p.ServiceType)
	d.Set("create_time", p.CreateTime)
	d.Set("created_by", p.CreatedBy)
	d.Set("update_time", p.UpdateTime)
	d.Set("updated_by", p.UpdatedBy)
	d.Set("guid", p.Guid)

	d.SetId(d.Id())
	return err
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	var err error
	client := m.(*policy.PolicyClient)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.DeletePolicyById(id)
	if err != nil {
		return err
	}

	return err
}
