---
subcategory: "Fault Injection Testing"
description: |-
  Fetches the details of a Fault Injection Testing Exclusion Window.
---

# google_fault_injection_testing_exclusion_window

Use this data source to get information about a Fault Injection Testing Exclusion Window. For more details, see the [API documentation](https://cloud.google.com/fault-injection/docs/reference/rest/v1alpha/projects.locations.exclusionWindows).

## Example Usage

```hcl
data "google_fault_injection_testing_exclusion_window" "default" {
  exclusion_window_id = "my-exclusion-window"
  location            = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `exclusion_window_id` -
  (Required)
  The ID of the Exclusion Window.

* `location` -
  (Required)
  The location in which the Exclusion Window resides.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_fault_injection_testing_exclusion_window](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/fault_injection_testing_exclusion_window) resource for details of all the available attributes.
