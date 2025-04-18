---
subcategory: "Looker"
description: |-
  Fetches the details of a Google Cloud Looker instance.
---

# google_looker_instance

Use this data source to get information about a Google Cloud Looker instance. For more details, see the [API documentation](https://cloud.google.com/looker/docs/reference/rest/v1/projects.locations.instances).

## Example Usage

```hcl
data "google_looker_instance" "default" {
  name   = "my-instance"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The ID of the instance or a fully qualified identifier for the instance. The ID must be 1-63 characters long, and comply with RFC 1035.

* `region` -
  (Required)
  The name of the Looker region of the instance.

* `project` - 
  (optional) 
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_looker_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/looker_instance) resource for details of all the available attributes.

Note that the `oauth_config` block is not available in the data source as it contains sensitive information.
