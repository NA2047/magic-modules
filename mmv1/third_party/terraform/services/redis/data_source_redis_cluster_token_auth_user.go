package redis

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceRedisClusterTokenAuthUser() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisClusterTokenAuthUser().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "cluster", "user_id")

	extraSchema := map[string]*schema.Schema{
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `The name of the region of the Redis cluster.`,
		},
	}
	dsSchema = tpgresource.MergeSchemas(extraSchema, dsSchema)

	return &schema.Resource{
		Read:   dataSourceRedisClusterTokenAuthUserRead,
		Schema: dsSchema,
	}
}

func dataSourceRedisClusterTokenAuthUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	userId := tpgresource.GetResourceNameFromSelfLink(d.Get("user_id").(string))
	if err := d.Set("user_id", userId); err != nil {
		return fmt.Errorf("Error setting user_id: %s", err)
	}

	cluster := d.Get("cluster").(string)
	if parts := regexp.MustCompile(`^projects/([^/]+)/locations/([^/]+)/clusters/([^/]+)$`).FindStringSubmatch(cluster); parts != nil {
		if err := d.Set("project", parts[1]); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("region", parts[2]); err != nil {
			return fmt.Errorf("Error setting region: %s", err)
		}
	} else {
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}
		region, err := tpgresource.GetRegion(d, config)
		if err != nil {
			return err
		}
		clusterName := tpgresource.GetResourceNameFromSelfLink(cluster)
		cluster = fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, region, clusterName)
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("region", region); err != nil {
			return fmt.Errorf("Error setting region: %s", err)
		}
		if err := d.Set("cluster", cluster); err != nil {
			return fmt.Errorf("Error setting cluster: %s", err)
		}
	}

	id, err := tpgresource.ReplaceVars(d, config, "{{cluster}}/tokenAuthUsers/{{user_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceRedisClusterTokenAuthUserRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_redis_cluster_token_auth_user",
		ProductName: "redis",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceRedisClusterTokenAuthUser(),
	}.Register()
}
