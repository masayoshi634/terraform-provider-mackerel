package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	MackerelAPIKeyParamName  = "MACKEREL_API_KEY"
	MackerelBaseUrlParamName = "MACKEREL_BASE_URL"
	MackerelDefaultBaseUrl   = "https://api.mackerelio.com/"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(MackerelAPIKeyParamName, nil),
				Description: "your Mackerel APIKey",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(MackerelBaseUrlParamName, MackerelDefaultBaseUrl),
				Description: "Mackerel Base URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mackerel_host_monitor":       resourceMackerelHostMonitor(),
			"mackerel_service_monitor":    resourceMackerelServiceMonitor(),
			"mackerel_external_monitor":   resourceMackerelExternalMonitor(),
			"mackerel_expression_monitor": resourceMackerelExpressionMonitor(),
			"mackerel_dashboard":          resourceMackerelDashboard(),
			"mackerel_service":            resourceMackerelService(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ApiKey:  d.Get("api_key").(string),
		BaseURL: d.Get("base_url").(string),
	}

	return config.NewClient()
}
