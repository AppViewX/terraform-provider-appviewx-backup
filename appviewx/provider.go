package appviewx

import (
	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			constants.APPVIEWX_USERNAME: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_PASSWORD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_ENVIRONMENT_IP: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_ENVIRONMENT_PORT: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_ENVIRONMENT_Is_HTTPS: &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"appviewx_automation":  ResourceAutomationServer(),
			"appviewx_certificate": ResourceCertificateServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	appviewxEnvironment := config.AppViewXEnvironment{
		AppViewXUserName:        d.Get(constants.APPVIEWX_USERNAME).(string),
		AppViewXPassword:        d.Get(constants.APPVIEWX_PASSWORD).(string),
		AppViewXEnvironmentIP:   d.Get(constants.APPVIEWX_ENVIRONMENT_IP).(string),
		AppViewXEnvironmentPort: d.Get(constants.APPVIEWX_ENVIRONMENT_PORT).(string),
		AppViewXIsHTTPS:         d.Get(constants.APPVIEWX_ENVIRONMENT_Is_HTTPS).(bool),
	}
	return &appviewxEnvironment, nil
}
