// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parallelstore

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceParallelstoreInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceParallelstoreInstance().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "location")

	return &schema.Resource{
		Read:   dataSourceParallelstoreInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceParallelstoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Setting location field, as this is set as a required field in instance resource to build the url
	d.Set("location", location)

	err = resourceParallelstoreInstanceRead(d, meta)
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
