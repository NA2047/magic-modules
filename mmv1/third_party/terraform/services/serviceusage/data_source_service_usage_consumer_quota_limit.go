// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package serviceusage

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceServiceUsageConsumerQuotaLimit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceUsageConsumerQuotaLimitRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the project in which the resource belongs. If it is not provided, the provider project is used.",
			},
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service name for which the quota limit is being retrieved. For example: compute.googleapis.com",
			},
			"metric": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The metric name for which the quota limit is being retrieved. For example: compute.googleapis.com%2Fcpus",
			},
			"limit_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The limit name for which the quota limit is being retrieved. For example: %2Fproject%2Fregion",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name of the quota limit.",
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The limit unit. An example unit would be 1/{project}/{region}",
			},
			"is_precise": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this limit is precise or imprecise.",
			},
			"allows_admin_overrides": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether admin overrides are allowed on this limit",
			},
			"quota_buckets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Summary of the enforced quota buckets, organized by quota dimension, ordered from least specific to most specific.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effective_limit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The effective limit of this quota bucket. Equal to defaultLimit if there are no overrides.",
						},
						"default_limit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default limit of this quota bucket, as specified by the service configuration.",
						},
						"dimensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The dimensions of this quota bucket.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"producer_override": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Producer override on this quota bucket.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource name of the override.",
									},
									"override_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The overriding quota limit value.",
									},
									"dimensions": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "The dimensions of the override.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"consumer_override": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Consumer override on this quota bucket.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource name of the override.",
									},
									"override_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The overriding quota limit value.",
									},
									"dimensions": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "The dimensions of the override.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"admin_override": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Admin override on this quota bucket.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource name of the override.",
									},
									"override_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The overriding quota limit value.",
									},
									"dimensions": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "The dimensions of the override.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"producer_quota_policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Producer policy inherited from the closest ancestor of the current consumer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource name of the policy.",
									},
									"policy_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The quota policy value.",
									},
									"dimensions": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "The dimensions of the policy.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"metric": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the metric to which this policy applies.",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The limit unit of the limit to which this policy applies.",
									},
									"container": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cloud resource container at which the quota policy is created.",
									},
								},
							},
						},
						"rollout_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rollout information of this quota bucket.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_limit_ongoing_rollout": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether there is an ongoing rollout for the default limit or not.",
									},
								},
							},
						},
					},
				},
			},
			"supported_locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all supported locations. This field is present only if the limit has a {region} or {zone} dimension.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceServiceUsageConsumerQuotaLimitRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	service := d.Get("service").(string)
	metric := d.Get("metric").(string)
	limitName := d.Get("limit_name").(string)

	url := fmt.Sprintf("%sprojects/%s/services/%s/consumerQuotaMetrics/%s/limits/%s", config.ServiceUsageBasePath, project, service, metric, limitName)

	log.Printf("[DEBUG] Reading ConsumerQuotaLimit %q", url)

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ServiceUsageConsumerQuotaLimit %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}

	if err := d.Set("name", flattenServiceUsageConsumerQuotaLimitName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}
	if err := d.Set("metric", flattenServiceUsageConsumerQuotaLimitMetric(res["metric"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}
	if err := d.Set("unit", flattenServiceUsageConsumerQuotaLimitUnit(res["unit"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}
	if err := d.Set("is_precise", flattenServiceUsageConsumerQuotaLimitIsPrecise(res["isPrecise"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}
	if err := d.Set("allows_admin_overrides", flattenServiceUsageConsumerQuotaLimitAllowsAdminOverrides(res["allowsAdminOverrides"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}
	if err := d.Set("quota_buckets", flattenServiceUsageConsumerQuotaLimitQuotaBuckets(res["quotaBuckets"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}
	if err := d.Set("supported_locations", flattenServiceUsageConsumerQuotaLimitSupportedLocations(res["supportedLocations"], d, config)); err != nil {
		return fmt.Errorf("Error reading ConsumerQuotaLimit: %s", err)
	}

	// Set the ID
	id := fmt.Sprintf("projects/%s/services/%s/consumerQuotaMetrics/%s/limits/%s", project, service, metric, limitName)
	d.SetId(id)

	return nil
}

func flattenServiceUsageConsumerQuotaLimitName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitMetric(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitUnit(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitIsPrecise(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitAllowsAdminOverrides(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitSupportedLocations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBuckets(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"effective_limit":       flattenServiceUsageConsumerQuotaLimitQuotaBucketsEffectiveLimit(original["effectiveLimit"], d, config),
			"default_limit":         flattenServiceUsageConsumerQuotaLimitQuotaBucketsDefaultLimit(original["defaultLimit"], d, config),
			"dimensions":            flattenServiceUsageConsumerQuotaLimitQuotaBucketsDimensions(original["dimensions"], d, config),
			"producer_override":     flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerOverride(original["producerOverride"], d, config),
			"consumer_override":     flattenServiceUsageConsumerQuotaLimitQuotaBucketsConsumerOverride(original["consumerOverride"], d, config),
			"admin_override":        flattenServiceUsageConsumerQuotaLimitQuotaBucketsAdminOverride(original["adminOverride"], d, config),
			"producer_quota_policy": flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicy(original["producerQuotaPolicy"], d, config),
			"rollout_info":          flattenServiceUsageConsumerQuotaLimitQuotaBucketsRolloutInfo(original["rolloutInfo"], d, config),
		})
	}
	return transformed
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsEffectiveLimit(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsDefaultLimit(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsDimensions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerOverride(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := []interface{}{
		map[string]interface{}{
			"name":           flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideName(original["name"], d, config),
			"override_value": flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideValue(original["overrideValue"], d, config),
			"dimensions":     flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideDimensions(original["dimensions"], d, config),
		},
	}
	return transformed
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsConsumerOverride(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := []interface{}{
		map[string]interface{}{
			"name":           flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideName(original["name"], d, config),
			"override_value": flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideValue(original["overrideValue"], d, config),
			"dimensions":     flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideDimensions(original["dimensions"], d, config),
		},
	}
	return transformed
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsAdminOverride(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := []interface{}{
		map[string]interface{}{
			"name":           flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideName(original["name"], d, config),
			"override_value": flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideValue(original["overrideValue"], d, config),
			"dimensions":     flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideDimensions(original["dimensions"], d, config),
		},
	}
	return transformed
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideValue(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsOverrideDimensions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := []interface{}{
		map[string]interface{}{
			"name":         flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyName(original["name"], d, config),
			"policy_value": flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyValue(original["policyValue"], d, config),
			"dimensions":   flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyDimensions(original["dimensions"], d, config),
			"metric":       flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyMetric(original["metric"], d, config),
			"unit":         flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyUnit(original["unit"], d, config),
			"container":    flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyContainer(original["container"], d, config),
		},
	}
	return transformed
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyValue(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyDimensions(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyMetric(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyUnit(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsProducerQuotaPolicyContainer(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsRolloutInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := []interface{}{
		map[string]interface{}{
			"default_limit_ongoing_rollout": flattenServiceUsageConsumerQuotaLimitQuotaBucketsRolloutInfoDefaultLimitOngoingRollout(original["defaultLimitOngoingRollout"], d, config),
		},
	}
	return transformed
}

func flattenServiceUsageConsumerQuotaLimitQuotaBucketsRolloutInfoDefaultLimitOngoingRollout(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
