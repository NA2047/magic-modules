package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeAcceleratorTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGoogleComputeAcceleratorTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A filter expression that filters resources listed in the response.`,
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The zone to list accelerator types for. If not specified, accelerator types from all zones will be returned.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs.`,
			},
			"types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of accelerator types.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the accelerator type.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `An optional textual description of the resource.`,
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the zone where the accelerator type resides.`,
						},
						"maximum_cards_per_instance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Maximum number of accelerator cards allowed per instance.`,
						},
						"self_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server-defined URL for the resource.`,
						},
						"creation_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Creation timestamp in RFC3339 text format.`,
						},
						"deprecated": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: `The deprecation status associated with this accelerator type.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The deprecation state of this resource. This can be ACTIVE, DEPRECATED, OBSOLETE, or DELETED.`,
									},
									"replacement": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The URL of the suggested replacement for a deprecated resource.`,
									},
									"deprecated": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `An optional RFC3339 timestamp on or after which the state of this resource is intended to change to DEPRECATED.`,
									},
									"obsolete": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `An optional RFC3339 timestamp on or after which the state of this resource is intended to change to OBSOLETE.`,
									},
									"deleted": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `An optional RFC3339 timestamp on or after which the state of this resource is intended to change to DELETED.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeAcceleratorTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	filter := d.Get("filter").(string)
	zone := d.Get("zone").(string)

	acceleratorTypes := make([]map[string]interface{}, 0)

	if zone != "" {
		// List accelerator types for a specific zone
		types, err := listAcceleratorTypesForZone(ctx, config, userAgent, project, zone, filter)
		if err != nil {
			return diag.FromErr(err)
		}
		acceleratorTypes = append(acceleratorTypes, types...)
	} else {
		// List accelerator types for all zones using aggregated list
		types, err := listAcceleratorTypesAggregated(ctx, config, userAgent, project, filter)
		if err != nil {
			return diag.FromErr(err)
		}
		acceleratorTypes = append(acceleratorTypes, types...)
	}

	if err := d.Set("types", acceleratorTypes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting types: %w", err))
	}

	if err := d.Set("project", project); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project: %w", err))
	}

	if err := d.Set("zone", zone); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone: %w", err))
	}

	id := fmt.Sprintf("projects/%s/acceleratorTypes", project)
	if zone != "" {
		id = fmt.Sprintf("projects/%s/zones/%s/acceleratorTypes", project, zone)
	}
	d.SetId(id)

	return nil
}

func listAcceleratorTypesForZone(ctx context.Context, config *transport_tpg.Config, userAgent, project, zone, filter string) ([]map[string]interface{}, error) {
	acceleratorTypes := make([]map[string]interface{}, 0)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewComputeClient(userAgent).AcceleratorTypes.List(project, zone).Context(ctx).Filter(filter).PageToken(token).Do()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving accelerator types for zone %s: %w", zone, err)
		}

		pageAcceleratorTypes := flattenAcceleratorTypesList(resp.Items)
		acceleratorTypes = append(acceleratorTypes, pageAcceleratorTypes...)

		token = resp.NextPageToken
		paginate = token != ""
	}

	return acceleratorTypes, nil
}

func listAcceleratorTypesAggregated(ctx context.Context, config *transport_tpg.Config, userAgent, project, filter string) ([]map[string]interface{}, error) {
	acceleratorTypes := make([]map[string]interface{}, 0)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewComputeClient(userAgent).AcceleratorTypes.AggregatedList(project).Context(ctx).Filter(filter).PageToken(token).Do()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving aggregated accelerator types: %w", err)
		}

		for _, scopedList := range resp.Items {
			if scopedList.AcceleratorTypes != nil {
				pageAcceleratorTypes := flattenAcceleratorTypesList(scopedList.AcceleratorTypes)
				acceleratorTypes = append(acceleratorTypes, pageAcceleratorTypes...)
			}
		}

		token = resp.NextPageToken
		paginate = token != ""
	}

	return acceleratorTypes, nil
}

func flattenAcceleratorTypesList(v []*compute.AcceleratorType) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	acceleratorTypes := make([]map[string]interface{}, 0, len(v))
	for _, at := range v {
		acceleratorType := map[string]interface{}{
			"name":                       at.Name,
			"description":                at.Description,
			"zone":                       at.Zone,
			"maximum_cards_per_instance": at.MaximumCardsPerInstance,
			"self_link":                  at.SelfLink,
			"creation_timestamp":         at.CreationTimestamp,
		}

		if dep := at.Deprecated; dep != nil {
			d := map[string]interface{}{
				"state":       dep.State,
				"replacement": dep.Replacement,
				"deprecated":  dep.Deprecated,
				"obsolete":    dep.Obsolete,
				"deleted":     dep.Deleted,
			}
			acceleratorType["deprecated"] = []map[string]interface{}{d}
		}

		acceleratorTypes = append(acceleratorTypes, acceleratorType)
	}

	return acceleratorTypes
}
