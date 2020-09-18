package provider

import (
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	ApiKey  string
	BaseURL string
}

func (c *Config) NewClient() (*mackerel.Client, error) {
	client, err := mackerel.NewClientWithOptions(c.ApiKey, c.BaseURL, false)
	if err != nil {
		return nil, err
	}

	client.UserAgent = "Terraform for Mackerel"
	if logging.IsDebugOrHigher() {
		client.Verbose = true
	}
	return client, nil
}
