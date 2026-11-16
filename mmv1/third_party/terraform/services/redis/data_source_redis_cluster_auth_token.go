package redis

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceRedisClusterAuthToken() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisClusterAuthToken().Schema)
	dsSchema["token"].Sensitive = true

	// Set 'Required' and 'Optional' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "cluster", "user_id", "token_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:          dataSourceRedisClusterAuthTokenRead,
		Schema:        dsSchema,
		UseJSONNumber: true,
	}
}

func dataSourceRedisClusterAuthTokenRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	cluster := d.Get("cluster").(string)
	if parts := redisClusterFullIdRegex.FindStringSubmatch(cluster); parts != nil {
		if d.Get("project") == "" {
			if err := d.Set("project", parts[1]); err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
		}
		if d.Get("region") == "" {
			if err := d.Set("region", parts[2]); err != nil {
				return fmt.Errorf("Error setting region: %s", err)
			}
		}
		if err := d.Set("cluster", parts[3]); err != nil {
			return fmt.Errorf("Error setting cluster: %s", err)
		}
	} else {
		if err := d.Set("cluster", tpgresource.GetResourceNameFromSelfLink(cluster)); err != nil {
			return fmt.Errorf("Error setting cluster: %s", err)
		}
	}

	userId := d.Get("user_id").(string)
	if parts := redisClusterUserFullIdRegex.FindStringSubmatch(userId); parts != nil {
		if d.Get("project") == "" {
			if err := d.Set("project", parts[1]); err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
		}
		if d.Get("region") == "" {
			if err := d.Set("region", parts[2]); err != nil {
				return fmt.Errorf("Error setting region: %s", err)
			}
		}
		if d.Get("cluster") == "" {
			if err := d.Set("cluster", parts[3]); err != nil {
				return fmt.Errorf("Error setting cluster: %s", err)
			}
		}
		if err := d.Set("user_id", parts[4]); err != nil {
			return fmt.Errorf("Error setting user_id: %s", err)
		}
	} else {
		if err := d.Set("user_id", tpgresource.GetResourceNameFromSelfLink(userId)); err != nil {
			return fmt.Errorf("Error setting user_id: %s", err)
		}
	}

	tokenId := d.Get("token_id").(string)
	if parts := redisClusterTokenFullIdRegex.FindStringSubmatch(tokenId); parts != nil {
		if d.Get("project") == "" {
			if err := d.Set("project", parts[1]); err != nil {
				return fmt.Errorf("Error setting project: %s", err)
			}
		}
		if d.Get("region") == "" {
			if err := d.Set("region", parts[2]); err != nil {
				return fmt.Errorf("Error setting region: %s", err)
			}
		}
		if d.Get("cluster") == "" {
			if err := d.Set("cluster", parts[3]); err != nil {
				return fmt.Errorf("Error setting cluster: %s", err)
			}
		}
		if d.Get("user_id") == "" {
			if err := d.Set("user_id", parts[4]); err != nil {
				return fmt.Errorf("Error setting user_id: %s", err)
			}
		}
		if err := d.Set("token_id", parts[5]); err != nil {
			return fmt.Errorf("Error setting token_id: %s", err)
		}
	} else {
		if err := d.Set("token_id", tpgresource.GetResourceNameFromSelfLink(tokenId)); err != nil {
			return fmt.Errorf("Error setting token_id: %s", err)
		}
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{region}}/clusters/{{cluster}}/tokenAuthUsers/{{user_id}}/authTokens/{{token_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	if err := resourceRedisClusterAuthTokenRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	if err := d.Set("deletion_policy", nil); err != nil {
		return fmt.Errorf("Error setting deletion_policy: %s", err)
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_redis_cluster_auth_token",
		ProductName: "redis",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceRedisClusterAuthToken(),
	}.Register()
}
