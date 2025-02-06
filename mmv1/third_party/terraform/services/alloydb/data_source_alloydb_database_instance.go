package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbDatabaseInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceAlloydbInstance().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "cluster", "instance_id")

	return &schema.Resource{
		Read:   dataSourceAlloydbDatabaseInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceAlloydbDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// location, err := tpgresource.GetLocation(d, config)
	// if err != nil {
	// 	return err
	// }

	// project, err := tpgresource.GetProject(d, config)
	// if err != nil {
	// 	return err
	// }	d.SetId(fmt.Sprintf("projects/%s/locations", project))
	GET https://alloydb.googleapis.com/v1/{parent=projects/*/locations/*}/supportedDatabaseFlags
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/supportedDbFlags", project, location))

	GET https://alloydb.googleapis.com/v1/{name=projects/*/locations/*}
	d.SetId(fmt.Sprintf("projects/%s/locations", project))

	GET https://alloydb.googleapis.com/v1/{name=projects/*/locations/*/clusters/*/instances/*}
	projects/{project}/locations/{region}/clusters/{clusterId}
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/clusters/%s/instances/%s/", project, location))

	id := fmt.Sprintf("projects/%s/instances/%s", project, location, d.Get("cluster").(string))
	if err != nil {
		return err
	}
	d.SetId(id)

	// id, err := tpgresource.ReplaceVars(d, config, "{{cluster}}/instances/{{instance_id}}")
	// if err != nil {
	// 	return fmt.Errorf("Error constructing id: %s", err)
	// }
	// d.SetId(id)

	err = resourceAlloydbInstanceRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
