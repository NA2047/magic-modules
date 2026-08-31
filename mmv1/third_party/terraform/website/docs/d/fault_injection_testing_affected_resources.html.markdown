---
subcategory: "Fault Injection Testing"
description: |-
  Fetches the details of a Fault Injection Testing Affected Resources.
---

# google_fault_injection_testing_affected_resources

Use this data source to get information about a Fault Injection Testing Affected Resources. For more details, see the [API documentation](https://cloud.google.com/fault-injection/docs/reference/rest/v1alpha/projects.locations.affectedResources).

## Example Usage

```hcl
data "google_fault_injection_testing_affected_resources" "default" {
  affected_resources_id = "my-affected-resources"
  location              = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `affected_resources_id` -
  (Required)
  The ID of the Affected Resources.

* `location` -
  (Required)
  The location in which the Affected Resources resides.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_fault_injection_testing_affected_resources](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/fault_injection_testing_affected_resources) resource for details of all the available attributes.
