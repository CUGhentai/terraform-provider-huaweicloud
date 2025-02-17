package servicestagev3

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	applicationsv2 "github.com/chnsz/golangsdk/openstack/servicestage/v2/applications"
	"github.com/chnsz/golangsdk/openstack/servicestage/v3/applications"
	"github.com/chnsz/golangsdk/openstack/servicestage/v3/components"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]([\w-]*[A-Za-z0-9])?$`),
						"The name must start with a letter and end with a letter or digit, and can only contain "+
							"letters, digits, underscores (_) and hyphens (-)."),
					validation.StringLenBetween(2, 64),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"environment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"variable": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_][\w-.]*$`),
												"The name can only contain letters, digits, underscores (_), "+
													"hyphens (-) and dots (.), and cannot start with a digit or dot."),
											validation.StringLenBetween(1, 64),
										),
									},
									"value": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 2048),
									},
								},
							},
						},
					},
				},
			},
			"component_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func createEnvironmentVariables(clientV2 *golangsdk.ServiceClient, appId string, envSet *schema.Set) error {
	for _, env := range envSet.List() {
		envMap := env.(map[string]interface{})
		varSet := envMap["variable"].(*schema.Set)
		variables := make([]applicationsv2.Variable, varSet.Len())
		for i, v := range varSet.List() {
			variable := v.(map[string]interface{})
			variables[i] = applicationsv2.Variable{
				Name:  variable["name"].(string),
				Value: variable["value"].(string),
			}
		}

		envId := envMap["id"].(string)
		appConfig := applicationsv2.Configuration{
			EnvVariables: variables,
		}
		if _, err := applicationsv2.UpdateConfig(clientV2, appId, envId, appConfig); err != nil {
			return err
		}
	}
	return nil
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	clientV2, err := hwConfig.ServiceStageV2Client(hwConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	client, err := hwConfig.ServiceStageV3Client(hwConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage v3 client: %s", err)
	}

	desc := d.Get("description").(string)
	opt := applications.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         &desc,
		EnterpriseProjectId: hwConfig.GetEnterpriseProjectID(d),
	}
	log.Printf("[DEBUG] The CreateOpts of ServiceStage application is: %v", opt)
	resp, err := applications.Create(client, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage application: %s", err)
	}

	d.SetId(resp.ID)

	err = createEnvironmentVariables(clientV2, d.Id(), d.Get("environment").(*schema.Set))
	if err != nil {
		return diag.Errorf("error creating environment variable: %s", err)
	}

	return resourceApplicationRead(ctx, d, meta)
}

func flattenEnvironmentVariables(configs []applicationsv2.ConfigResp) []map[string]interface{} {
	if len(configs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(configs))

	for i, v := range configs {
		variables := make([]map[string]interface{}, len(v.Configuration.EnvVariables))
		for j, variable := range v.Configuration.EnvVariables {
			variables[j] = map[string]interface{}{
				"name":  variable.Name,
				"value": variable.Value,
			}
		}

		result[i] = map[string]interface{}{
			"id":       v.EnvironmentId,
			"variable": variables,
		}
	}
	return result
}

func flattenComponentIds(list []components.Component) []string {
	result := make([]string, len(list))

	for i, component := range list {
		result[i] = component.ID
	}
	return result
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	region := hwConfig.GetRegion(d)

	clientV2, err := hwConfig.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	client, err := hwConfig.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v3 client: %s", err)
	}

	application, err := applications.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage application")
	}

	componentList, err := components.List(client, d.Id(), components.ListOpts{})
	if err != nil {
		return diag.Errorf("error getting components under application (%s): %s", d.Id(), err)
	}

	configs, err := applicationsv2.ListConfig(clientV2, d.Id(), applicationsv2.ListConfigOpts{})
	if err != nil {
		return diag.Errorf("error getting environment variables under application (%s): %s", d.Id(), err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", application.Name),
		d.Set("description", application.Description),
		d.Set("enterprise_project_id", application.EnterpriseProjectId),
		d.Set("environment", flattenEnvironmentVariables(configs)),
		d.Set("component_ids", flattenComponentIds(componentList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func removeEnvironmentVariables(clientV2 *golangsdk.ServiceClient, appId string, envSet *schema.Set) error {
	for _, env := range envSet.List() {
		envMap := env.(map[string]interface{})
		envId := envMap["id"].(string)
		if err := applicationsv2.DeleteConfig(clientV2, appId, envId).ExtractErr(); err != nil {
			return err
		}
	}
	return nil
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	region := hwConfig.GetRegion(d)

	clientV2, err := hwConfig.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		desc := d.Get("description").(string)
		updateOpt := applicationsv2.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: &desc,
		}
		_, err = applicationsv2.Update(clientV2, d.Id(), updateOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceStage application (%s): %s", d.Id(), err)
		}
	}

	if d.HasChanges("environment") {
		oldRaws, newRaws := d.GetChange("environment")
		addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
		if err := removeEnvironmentVariables(clientV2, d.Id(), removeRaws); err != nil {
			return diag.FromErr(err)
		}
		if err := createEnvironmentVariables(clientV2, d.Id(), addRaws); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hwConfig := meta.(*config.Config)
	region := hwConfig.GetRegion(d)
	client, err := hwConfig.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v3 client: %s", err)
	}

	err = applications.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting ServiceStage application (%s): %s", d.Id(), err)
	}
	return nil
}
