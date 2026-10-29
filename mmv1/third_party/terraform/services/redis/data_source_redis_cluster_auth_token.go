package redis

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceRedisClusterAuthToken() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceRedisClusterAuthToken().Schema)

	// Set 'Required' and 'Optional' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "token_auth_user")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "token_id")

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
		Read:   dataSourceRedisClusterAuthTokenRead,
		Schema: dsSchema,
	}
}

func dataSourceRedisClusterAuthTokenRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	tokenAuthUser := d.Get("token_auth_user").(string)
	if parts := regexp.MustCompile(`^projects/([^/]+)/locations/([^/]+)/clusters/([^/]+)/tokenAuthUsers/([^/]+)$`).FindStringSubmatch(tokenAuthUser); parts != nil {
		if err := d.Set("project", parts[1]); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("region", parts[2]); err != nil {
			return fmt.Errorf("Error setting region: %s", err)
		}
	} else if parts := regexp.MustCompile(`^(?:clusters/)?([^/]+)/tokenAuthUsers/([^/]+)$`).FindStringSubmatch(tokenAuthUser); parts != nil {
		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}
		region, err := tpgresource.GetRegion(d, config)
		if err != nil {
			return err
		}
		cluster := parts[1]
		userId := parts[2]
		tokenAuthUser = fmt.Sprintf("projects/%s/locations/%s/clusters/%s/tokenAuthUsers/%s", project, region, cluster, userId)
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("region", region); err != nil {
			return fmt.Errorf("Error setting region: %s", err)
		}
		if err := d.Set("token_auth_user", tokenAuthUser); err != nil {
			return fmt.Errorf("Error setting token_auth_user: %s", err)
		}
	} else {
		return fmt.Errorf("Invalid token_auth_user format: %q. Expected projects/{project}/locations/{region}/clusters/{cluster}/tokenAuthUsers/{user_id} or {cluster}/tokenAuthUsers/{user_id}", tokenAuthUser)
	}

	if v, ok := d.GetOk("token_id"); ok && v.(string) != "" {
		tokenId := tpgresource.GetResourceNameFromSelfLink(v.(string))
		if err := d.Set("token_id", tokenId); err != nil {
			return fmt.Errorf("Error setting token_id: %s", err)
		}
		id, err := tpgresource.ReplaceVars(d, config, "{{token_auth_user}}/authTokens/{{token_id}}")
		if err != nil {
			return fmt.Errorf("Error constructing id: %s", err)
		}
		d.SetId(id)

		err = resourceRedisClusterAuthTokenRead(d, meta)
		if err != nil {
			return err
		}
		if d.Id() == "" {
			return fmt.Errorf("%s not found", id)
		}
		return nil
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{RedisBasePath}}{{token_auth_user}}/authTokens")
	if err != nil {
		return err
	}

	billingProject := ""
	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}
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
		return fmt.Errorf("Error listing authTokens for %s: %s", url, err)
	}

	tokens, ok := res["authTokens"].([]interface{})
	if !ok || len(tokens) == 0 {
		return fmt.Errorf("No authTokens found for %s", d.Get("token_auth_user"))
	}

	var latestToken map[string]interface{}
	maxVersion := -1
	for _, t := range tokens {
		if tm, ok := t.(map[string]interface{}); ok {
			if name, ok := tm["name"].(string); ok {
				tokenId := tpgresource.GetResourceNameFromSelfLink(name)
				var vNum int
				if _, err := fmt.Sscanf(tokenId, "version-%d", &vNum); err == nil {
					if vNum > maxVersion {
						maxVersion = vNum
						latestToken = tm
					}
				} else if latestToken == nil {
					latestToken = tm
				}
			}
		}
	}
	if latestToken == nil {
		latestToken = tokens[len(tokens)-1].(map[string]interface{})
	}

	tokenName, ok := latestToken["name"].(string)
	if !ok {
		return fmt.Errorf("AuthToken response missing name attribute")
	}
	tokenId := tpgresource.GetResourceNameFromSelfLink(tokenName)
	d.SetId(tokenName)
	if err := d.Set("token_id", tokenId); err != nil {
		return fmt.Errorf("Error setting token_id: %s", err)
	}
	if err := d.Set("name", tokenName); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	err = resourceRedisClusterAuthTokenRead(d, meta)
	if err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("%s not found", tokenName)
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
