---
subcategory: "Fault Injection Testing"
description: |-
  Fetches the details of a Fault Injection Testing Validation.
---

# google_fault_injection_testing_validation

Use this data source to get information about a Fault Injection Testing Validation. For more details, see the [API documentation](https://cloud.google.com/fault-injection/docs/reference/rest/v1alpha/projects.locations.validations).

## Example Usage

```hcl
data "google_fault_injection_testing_validation" "default" {
  validation_id = "my-validation"
  location      = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `validation_id` -
  (Required)
  The ID of the Validation.

* `location` -
  (Required)
  The location in which the Validation resides.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_fault_injection_testing_validation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/fault_injection_testing_validation) resource for details of all the available attributes.
